/*
 * // Copyright (c) 2017 Cisco and/or its affiliates.
 * //
 * // Licensed under the Apache License, Version 2.0 (the "License");
 * // you may not use this file except in compliance with the License.
 * // You may obtain a copy of the License at:
 * //
 * //     http://www.apache.org/licenses/LICENSE-2.0
 * //
 * // Unless required by applicable law or agreed to in writing, software
 * // distributed under the License is distributed on an "AS IS" BASIS,
 * // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * // See the License for the specific language governing permissions and
 * // limitations under the License.
 */

package configurator

import (
	"net"
	"sort"

	"go.ligato.io/cn-infra/v2/logging"

	podmodel "github.com/americanbinary/vpp/plugins/ksr/model/pod"
	"github.com/americanbinary/vpp/plugins/policy/cache"
	"github.com/americanbinary/vpp/plugins/policy/renderer"
	"github.com/americanbinary/vpp/plugins/policy/utils"
)

// PolicyConfigurator translates a set of Contiv Policies into ingress and
// egress lists of Contiv Rules (n-tuples with the most basic policy rule
// definition) and applies them into the target vswitch via registered
// renderers. Allows to register multiple renderers for different network stacks.
// For the best performance, creates a shortest possible sequence of rules
// that implement a given policy. Furthermore, to allow renderers share a list
// of ingress or egress rules between interfaces, the same set of policies
// always results in the same list of rules.
type PolicyConfigurator struct {
	Deps

	renderers         []renderer.PolicyRendererAPI
	parallelRendering bool
	podIPAddresses    PodIPAddresses
}

// Deps lists dependencies of PolicyConfigurator.
type Deps struct {
	Log   logging.Logger
	Cache cache.PolicyCacheAPI
	IPAM  IPAM
}

// IPAM interface lists IPAM methods needed by Policy Configurator.
type IPAM interface {
	// NatLoopbackIP returns the IP address of a virtual loopback, used to route
	// traffic between clients and services via VPP even if the source and destination
	// are the same IP addresses and would otherwise be routed locally.
	NatLoopbackIP() net.IP
}

// PolicyConfiguratorTxn represents a single transaction of the policy configurator.
type PolicyConfiguratorTxn struct {
	Log            logging.Logger
	configurator   *PolicyConfigurator
	resync         bool
	config         map[podmodel.ID]ContivPolicies // config to render
	podIPAddresses PodIPAddresses
}

// ContivPolicies is a list of policies that can be ordered by policy ID.
type ContivPolicies []*ContivPolicy

// ProcessedPolicySet stores configuration already generated for a given
// set of policies. It is used only temporarily for a duration of the commit
// for a performance optimization.
type ProcessedPolicySet struct {
	policies ContivPolicies // ordered
	ingress  *ContivRules
	egress   *ContivRules
}

// ContivRules is a list of Contiv rules without duplicities.
type ContivRules struct {
	rules        []*renderer.ContivRule // rules in the original order
	orderedRules []*renderer.ContivRule // ordered to test for duplicities
}

// PodIPAddresses is a map used to remember IP address for each configured pod.
type PodIPAddresses map[podmodel.ID]*net.IPNet

// Init initializes policy configurator.
func (pc *PolicyConfigurator) Init(parallelRendering bool) error {
	pc.renderers = []renderer.PolicyRendererAPI{}
	pc.parallelRendering = parallelRendering
	pc.podIPAddresses = make(PodIPAddresses)
	return nil
}

// RegisterRenderer registers a new renderer.
// The renderer will be receiving rules for all pods in this K8s node.
// It is up to the render to possibly filter out rules for pods without
// an inter-connection in the destination network stack.
func (pc *PolicyConfigurator) RegisterRenderer(renderer renderer.PolicyRendererAPI) error {
	pc.renderers = append(pc.renderers, renderer)
	return nil
}

// Close deallocates resource held by the configurator.
func (pc *PolicyConfigurator) Close() error {
	return nil
}

// NewTxn starts a new transaction. The re-configuration executes only after
// Commit() is called. If <resync> is enabled, the supplied configuration will
// completely replace the existing one, otherwise pods not mentioned in the
// transaction are left unchanged.
func (pc *PolicyConfigurator) NewTxn(resync bool) Txn {
	txn := &PolicyConfiguratorTxn{
		Log:          pc.Log,
		configurator: pc,
		resync:       resync,
		config:       make(map[podmodel.ID]ContivPolicies),
	}
	if resync {
		txn.podIPAddresses = make(PodIPAddresses)
	} else {
		txn.podIPAddresses = pc.podIPAddresses.Copy()
	}
	return txn
}

// Configure applies the set of policies for a given pod. The existing policies
// are replaced. The order of policies is not important (it is a set).
func (pct *PolicyConfiguratorTxn) Configure(pod podmodel.ID, policies []*ContivPolicy) Txn {
	pct.Log.WithFields(logging.Fields{
		"pod":      pod,
		"policies": policies,
	}).Debug("PolicyConfigurator Configure()")
	pct.config[pod] = policies
	return pct
}

// Commit proceeds with the reconfiguration.
func (pct *PolicyConfiguratorTxn) Commit() error {
	// Remember processed sets of policies between iterations so that the same
	// set will not be processed more than once.
	processed := []ProcessedPolicySet{}

	// Transactions of all registered renderers.
	rendererTxns := []renderer.Txn{}

	for pod, unorderedPolicies := range pct.config {
		ingress := &ContivRules{}
		egress := &ContivRules{}
		var delPodConfig bool

		// Get target pod configuration.
		podIPNet, hadIPAddr := pct.podIPAddresses[pod]
		found, podData := pct.configurator.Cache.LookupPod(pod)

		// Handle removed pod.
		if !found || podData.IpAddress == "" {
			if hadIPAddr {
				pct.Log.WithField("pod", pod).Debug("Removing policies from the pod.")
				delPodConfig = true
				delete(pct.podIPAddresses, pod)
			} else {
				/* already un-configured */
				continue
			}
		}

		if !delPodConfig {
			// Get pod IP address (expressed as one-host subnet).
			podIPNet = utils.GetOneHostSubnet(podData.IpAddress)
			if podIPNet == nil {
				pct.Log.WithField("pod", pod).Warn("Pod has invalid IP address assigned")
				continue
			}
			pct.podIPAddresses[pod] = podIPNet

			// Sort policies to get the same outcome for the same set.
			policies := unorderedPolicies.Copy()
			sort.Sort(policies)

			// Check if this set was already processed.
			alreadyProcessed := false
			for _, policySet := range processed {
				if policySet.policies.Equals(policies) {
					ingress = policySet.ingress
					egress = policySet.egress
					alreadyProcessed = true
				}
			}

			// Generate rules for a set of policies not yet processed.
			if !alreadyProcessed {
				// Direction in policies is from the pod point of view, whereas rules
				// are evaluated from the vswitch perspective.
				egress = pct.generateRules(MatchIngress, policies)
				ingress = pct.generateRules(MatchEgress, policies)
				// Remember already processed set of policies.
				processed = append(processed,
					ProcessedPolicySet{
						policies: policies,
						ingress:  ingress,
						egress:   egress,
					})
			}
		}

		// Start transaction on every renderer if they are not running already.
		if len(rendererTxns) == 0 {
			for _, renderer := range pct.configurator.renderers {
				rendererTxns = append(rendererTxns, renderer.NewTxn(pct.resync))
			}
		}

		// Add rules into the transactions.
		for _, rTxn := range rendererTxns {
			rTxn.Render(pod, podIPNet, ingress.CopySlice(), egress.CopySlice(), delPodConfig)
		}
	}

	// Commit all renderer transactions.
	var wasError error
	rndrChan := make(chan error)
	for _, rTxn := range rendererTxns {
		if pct.configurator.parallelRendering {
			go func(txn renderer.Txn) {
				err := txn.Commit()
				rndrChan <- err
			}(rTxn)
		} else {
			err := rTxn.Commit()
			if err != nil {
				wasError = err
			}
		}
	}
	if pct.configurator.parallelRendering {
		for i := 0; i < len(rendererTxns); i++ {
			err := <-rndrChan
			if err != nil {
				wasError = err
			}
		}
	}

	// Save changes to the configurator.
	pct.configurator.podIPAddresses = pct.podIPAddresses.Copy()

	return wasError
}

// PeerPod represents the opposite pod in the policy rule.
type PeerPod struct {
	ID    podmodel.ID
	IPNet *net.IPNet
}

// Generate a list of ingress or egress rules implementing a given list of policies.
func (pct *PolicyConfiguratorTxn) generateRules(direction MatchType, policies ContivPolicies) *ContivRules {
	rules := &ContivRules{}
	hasPolicy := false
	allAllowed := false

	for _, policy := range policies {
		if (policy.Type == PolicyIngress && direction == MatchEgress) ||
			(policy.Type == PolicyEgress && direction == MatchIngress) {
			// Policy does not apply to this direction.
			continue
		}
		hasPolicy = true

		for _, match := range policy.Matches {
			if match.Type != direction {
				continue
			}

			// Collect IP addresses of all pod peers.
			peers := []PeerPod{}
			for _, peer := range match.Pods {
				found, peerData := pct.configurator.Cache.LookupPod(peer)
				if !found {
					pct.Log.WithField("peer", peer).Warn("Peer pod data not found in the cache")
					continue
				}
				if peerData.IpAddress == "" {
					pct.Log.WithField("peer", peer).Debug("Peer pod has no IP address assigned")
					continue
				}
				peerIPNet := utils.GetOneHostSubnet(peerData.IpAddress)
				if peerIPNet == nil {
					pct.Log.WithFields(logging.Fields{
						"peer": peer,
						"ip":   peerData.IpAddress}).Warn("Peer pod has invalid IP address assigned")
					continue
				}
				peers = append(peers, PeerPod{ID: peer, IPNet: peerIPNet})
			}

			// Collect all subnets from IPBlocks.
			allSubnets := []*net.IPNet{}
			for _, block := range match.IPBlocks {
				subnets := []*net.IPNet{&block.Network}
				for _, except := range block.Except {
					subtracted := []*net.IPNet{}
					for _, subnet := range subnets {
						subtracted = append(subtracted, subtractSubnet(subnet, &except)...)
					}
					subnets = subtracted
				}
				allSubnets = append(allSubnets, subnets...)
			}

			// Handle undefined set of pods and IP blocks.
			// = match anything on L3
			if match.Pods == nil && match.IPBlocks == nil {
				if len(match.Ports) == 0 {
					// = match anything on L3 & L4
					ruleAny := &renderer.ContivRule{
						Action:      renderer.ActionPermit,
						SrcNetwork:  &net.IPNet{},
						DestNetwork: &net.IPNet{},
						Protocol:    renderer.ANY,
						SrcPort:     0,
						DestPort:    0,
					}
					rules.Insert(ruleAny)
					allAllowed = true
				} else {
					// = match by L4
					for _, port := range match.Ports {
						rule := &renderer.ContivRule{
							Action:      renderer.ActionPermit,
							SrcNetwork:  &net.IPNet{},
							DestNetwork: &net.IPNet{},
							SrcPort:     0,
							DestPort:    port.Number,
						}
						if port.Protocol == TCP {
							rule.Protocol = renderer.TCP
						} else {
							rule.Protocol = renderer.UDP
						}
						rules.Insert(rule)
					}
				}
			}

			// Combine pod peers with ports.
			for _, peer := range peers {
				if len(match.Ports) == 0 {
					// Match all ports.
					// = match by L3
					ruleAny := &renderer.ContivRule{
						Action:      renderer.ActionPermit,
						Protocol:    renderer.ANY,
						SrcNetwork:  &net.IPNet{},
						DestNetwork: &net.IPNet{},
						SrcPort:     0,
						DestPort:    0,
					}
					if direction == MatchIngress {
						ruleAny.SrcNetwork = peer.IPNet
					} else {
						ruleAny.DestNetwork = peer.IPNet
					}
					rules.Insert(ruleAny)
				} else {
					// Combine each port with the peer.
					// = match by L3 & L4
					for _, port := range match.Ports {
						rule := &renderer.ContivRule{
							Action:      renderer.ActionPermit,
							SrcNetwork:  &net.IPNet{},
							DestNetwork: &net.IPNet{},
							SrcPort:     0,
							DestPort:    port.Number,
						}
						if direction == MatchIngress {
							rule.SrcNetwork = peer.IPNet
						} else {
							rule.DestNetwork = peer.IPNet
						}
						if port.Protocol == TCP {
							rule.Protocol = renderer.TCP
						} else {
							rule.Protocol = renderer.UDP
						}
						rules.Insert(rule)
					}
				}
			}

			// Combine IPBlocks with ports.
			for _, subnet := range allSubnets {
				if len(match.Ports) == 0 {
					// Handle IPBlock with no ports.
					// = match by L3
					ruleAny := &renderer.ContivRule{
						Action:      renderer.ActionPermit,
						Protocol:    renderer.ANY,
						SrcNetwork:  &net.IPNet{},
						DestNetwork: &net.IPNet{},
						SrcPort:     0,
						DestPort:    0,
					}
					if direction == MatchIngress {
						ruleAny.SrcNetwork = subnet
					} else {
						ruleAny.DestNetwork = subnet
					}
					rules.Insert(ruleAny)
				} else {
					// Combine each port with the block.
					// = match by L3 & L4
					for _, port := range match.Ports {
						rule := &renderer.ContivRule{
							Action:      renderer.ActionPermit,
							SrcNetwork:  &net.IPNet{},
							DestNetwork: &net.IPNet{},
							SrcPort:     0,
							DestPort:    port.Number,
						}
						if direction == MatchIngress {
							rule.SrcNetwork = subnet
						} else {
							rule.DestNetwork = subnet
						}
						if port.Protocol == TCP {
							rule.Protocol = renderer.TCP
						} else {
							rule.Protocol = renderer.UDP
						}
						rules.Insert(rule)
					}
				}
			}
		}
	}

	if hasPolicy && !allAllowed {
		if direction == MatchIngress {
			// Allow connections from the virtual NAT-loopback (access to service from itself).
			natLoopIP := pct.configurator.IPAM.NatLoopbackIP()
			ruleAny := &renderer.ContivRule{
				Action:      renderer.ActionPermit,
				Protocol:    renderer.ANY,
				SrcNetwork:  utils.GetOneHostSubnetFromIP(natLoopIP),
				DestNetwork: &net.IPNet{},
				SrcPort:     0,
				DestPort:    0,
			}
			rules.Insert(ruleAny)
		}
		// Deny the rest.
		ruleNone := &renderer.ContivRule{
			Action:      renderer.ActionDeny,
			SrcNetwork:  &net.IPNet{},
			DestNetwork: &net.IPNet{},
			Protocol:    renderer.ANY,
			SrcPort:     0,
			DestPort:    0,
		}
		rules.Insert(ruleNone)
	}

	return rules
}

// Copy creates a shallow copy of ContivPolicies.
func (cp ContivPolicies) Copy() ContivPolicies {
	cpCopy := make(ContivPolicies, len(cp))
	copy(cpCopy, cp)
	return cpCopy
}

// Equals returns true for equal lists of policies.
func (cp ContivPolicies) Equals(cp2 ContivPolicies) bool {
	if len(cp) != len(cp2) {
		return false
	}
	for idx, policy := range cp {
		if policy.ID != cp2[idx].ID {
			return false
		}
	}
	return true
}

// Len return the number of policies in the list.
func (cp ContivPolicies) Len() int {
	return len(cp)
}

// Swap replaces order of two policies in the list.
func (cp ContivPolicies) Swap(i, j int) {
	cp[i], cp[j] = cp[j], cp[i]
}

// Less compares two policies by their IDs.
func (cp ContivPolicies) Less(i, j int) bool {
	if cp[i].ID.Namespace < cp[j].ID.Namespace {
		return true
	}
	if cp[i].ID.Namespace == cp[j].ID.Namespace {
		if cp[i].ID.Name < cp[j].ID.Name {
			return true
		}
	}
	return false
}

// Insert inserts the rule into the list.
// Returns *true* if the rule was inserted, *false* if the same rule is already
// in the list.
func (cr *ContivRules) Insert(rule *renderer.ContivRule) bool {
	// get the index at which the rule should be inserted to keep the order
	idx := sort.Search(len(cr.orderedRules),
		func(i int) bool {
			return rule.Compare(cr.orderedRules[i]) <= 0
		})
	if idx < len(cr.orderedRules) && rule.Compare(cr.orderedRules[idx]) == 0 {
		return false
	}

	// allocate new entry at the right index
	cr.orderedRules = append(cr.orderedRules, nil)
	if idx < len(cr.orderedRules) {
		copy(cr.orderedRules[idx+1:], cr.orderedRules[idx:])
	}

	// add entry into both internal lists
	cr.orderedRules[idx] = rule
	cr.rules = append(cr.rules, rule)
	return true
}

// CopySlice returns a deep-copied slice of all rules (in the order as inserted).
func (cr *ContivRules) CopySlice() []*renderer.ContivRule {
	slice := make([]*renderer.ContivRule, len(cr.rules))
	for idx, rule := range cr.rules {
		slice[idx] = &renderer.ContivRule{}
		*(slice[idx]) = *rule
	}
	return slice
}

// Copy creates a deep copy of PodIPAddresses.
func (pa PodIPAddresses) Copy() PodIPAddresses {
	paCopy := make(PodIPAddresses, len(pa))
	for pod, addr := range pa {
		paCopy[pod] = addr
	}
	return paCopy
}

// Function returns a list of subnets with all IPs included in net1 and not included in net2.
func subtractSubnet(net1, net2 *net.IPNet) []*net.IPNet {
	result := []*net.IPNet{}
	net1MaskSize, _ := net1.Mask.Size()
	net2MaskSize, _ := net2.Mask.Size()
	if net1MaskSize > net2MaskSize {
		// net2 higher than net1 in the tree
		if !net2.Contains(net1.IP) {
			result = append(result, net1)
		}
	} else if net1MaskSize == net2MaskSize {
		// same level in the tree
		if !net1.IP.Equal(net2.IP) {
			result = append(result, net1)
		}
	} else {
		// net2 lower then net1 in the tree
		if !net1.Contains(net2.IP) {
			result = append(result, net1)
		} else {
			// net2 under net1
			for bit := net1MaskSize; bit < net2MaskSize; bit++ {
				subnet := &net.IPNet{}
				subnet.Mask = net.CIDRMask(bit+1, len(net2.Mask)*8)
				subnet.IP = net2.IP.Mask(subnet.Mask)
				// flip the last bit of the IP
				subnet.IP[bit/8] ^= byte(1 << uint(7-(bit%8)))
				result = append(result, subnet)
			}
		}
	}

	return result
}
