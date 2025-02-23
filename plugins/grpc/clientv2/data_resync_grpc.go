// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clientv2

import (
	"golang.org/x/net/context"

	"github.com/americanbinary/vpp/plugins/grpc/rpc"

	"go.ligato.io/vpp-agent/v3/clientv2/linux"
	"go.ligato.io/vpp-agent/v3/clientv2/vpp"

	"go.ligato.io/vpp-agent/v3/proto/ligato/linux/interfaces"
	"go.ligato.io/vpp-agent/v3/proto/ligato/linux/iptables"
	"go.ligato.io/vpp-agent/v3/proto/ligato/linux/l3"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/abf"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/acl"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/ipsec"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l2"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l3"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/nat"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/punt"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/stn"
)

// NewDataResyncDSL is a constructor.
func NewDataResyncDSL(client rpc.DataResyncServiceClient) *DataResyncDSL {
	return &DataResyncDSL{
		client: client,
		req:    &rpc.DataRequest{},
	}
}

// DataResyncDSL is used to conveniently assign all the data that are needed for the RESYNC.
// This is implementation of Domain Specific Language (DSL) for data RESYNC of the VPP configuration.
type DataResyncDSL struct {
	client rpc.DataResyncServiceClient
	req    *rpc.DataRequest
}

// LinuxInterface adds Linux interface to the RESYNC request.
func (dsl *DataResyncDSL) LinuxInterface(val *linux_interfaces.Interface) linuxclient.DataResyncDSL {
	dsl.req.LinuxInterfaces = append(dsl.req.LinuxInterfaces, val)
	return dsl
}

// IptablesRuleChain adds iptables rule chain to the RESYNC request.
func (dsl *DataResyncDSL) IptablesRuleChain(val *linux_iptables.RuleChain) linuxclient.DataResyncDSL {
	dsl.req.RuleChains = append(dsl.req.RuleChains, val)
	return dsl
}

// Span adds VPP span to the RESYNC request.
func (dsl *DataResyncDSL) Span(span *vpp_interfaces.Span) linuxclient.DataResyncDSL {
	dsl.req.Spans = append(dsl.req.Spans, span)
	return dsl
}

// ABF adds ACL-based forwarding to the RESYNC request.
func (dsl *DataResyncDSL) ABF(abf *vpp_abf.ABF) linuxclient.DataResyncDSL {
	dsl.req.Forwardings = append(dsl.req.Forwardings, abf)
	return dsl
}

// VrfTable adds VPP VRF table to the RESYNC request.
func (dsl *DataResyncDSL) VrfTable(val *vpp_l3.VrfTable) linuxclient.DataResyncDSL {
	dsl.req.VrfTables = append(dsl.req.VrfTables, val)
	return dsl
}

// LinuxArpEntry adds Linux ARP entry to the RESYNC request.
func (dsl *DataResyncDSL) LinuxArpEntry(val *linux_l3.ARPEntry) linuxclient.DataResyncDSL {
	dsl.req.LinuxArpEntries = append(dsl.req.LinuxArpEntries, val)
	return dsl
}

// NAT44Interface adds NAT44 interface configuration to the RESYNC request.
func (dsl *DataResyncDSL) NAT44Interface(natIf *vpp_nat.Nat44Interface) linuxclient.DataResyncDSL {
	dsl.req.Nat44Interfaces = append(dsl.req.Nat44Interfaces, natIf)
	return dsl
}

// NAT44AddressPool adds NAT44 address pool configuration to the RESYNC request.
func (dsl *DataResyncDSL) NAT44AddressPool(pool *vpp_nat.Nat44AddressPool) linuxclient.DataResyncDSL {
	dsl.req.Nat44AddressPools = append(dsl.req.Nat44AddressPools, pool)
	return dsl
}

// PuntException adds request to create or update exception to punt specific packets.
func (dsl *DataResyncDSL) PuntException(val *vpp_punt.Exception) linuxclient.DataResyncDSL {
	dsl.req.Exceptions = append(dsl.req.Exceptions, val)
	return dsl
}

// LinuxRoute adds Linux route to the RESYNC request.
func (dsl *DataResyncDSL) LinuxRoute(val *linux_l3.Route) linuxclient.DataResyncDSL {
	dsl.req.LinuxRoutes = append(dsl.req.LinuxRoutes, val)
	return dsl
}

// VppInterface adds VPP interface to the RESYNC request.
func (dsl *DataResyncDSL) VppInterface(val *vpp_interfaces.Interface) linuxclient.DataResyncDSL {
	dsl.req.Interfaces = append(dsl.req.Interfaces, val)
	return dsl
}

// ACL adds VPP Access Control List to the RESYNC request.
func (dsl *DataResyncDSL) ACL(val *vpp_acl.ACL) linuxclient.DataResyncDSL {
	dsl.req.AccessLists = append(dsl.req.AccessLists, val)
	return dsl
}

// BD adds VPP Bridge Domain to the RESYNC request.
func (dsl *DataResyncDSL) BD(val *vpp_l2.BridgeDomain) linuxclient.DataResyncDSL {
	dsl.req.BridgeDomains = append(dsl.req.BridgeDomains, val)
	return dsl
}

// BDFIB adds VPP L2 FIB to the RESYNC request.
func (dsl *DataResyncDSL) BDFIB(val *vpp_l2.FIBEntry) linuxclient.DataResyncDSL {
	dsl.req.FIBs = append(dsl.req.FIBs, val)
	return dsl
}

// XConnect adds VPP Cross Connect to the RESYNC request.
func (dsl *DataResyncDSL) XConnect(val *vpp_l2.XConnectPair) linuxclient.DataResyncDSL {
	dsl.req.XCons = append(dsl.req.XCons, val)
	return dsl
}

// StaticRoute adds VPP L3 Static Route to the RESYNC request.
func (dsl *DataResyncDSL) StaticRoute(val *vpp_l3.Route) linuxclient.DataResyncDSL {
	dsl.req.StaticRoutes = append(dsl.req.StaticRoutes, val)
	return dsl
}

// Arp adds VPP L3 ARP to the RESYNC request.
func (dsl *DataResyncDSL) Arp(val *vpp_l3.ARPEntry) linuxclient.DataResyncDSL {
	dsl.req.ArpEntries = append(dsl.req.ArpEntries, val)
	return dsl
}

// ProxyArp adds L3 proxy ARP interfaces to the RESYNC request.
func (dsl *DataResyncDSL) ProxyArp(val *vpp_l3.ProxyARP) linuxclient.DataResyncDSL {
	dsl.req.ProxyArp = val
	return dsl
}

// IPScanNeighbor adds L3 IP Scan Neighbor to the RESYNC request.
func (dsl *DataResyncDSL) IPScanNeighbor(val *vpp_l3.IPScanNeighbor) linuxclient.DataResyncDSL {
	dsl.req.IPScanNeighbor = val
	return dsl
}

// StnRule adds Stn rule to the RESYNC request.
func (dsl *DataResyncDSL) StnRule(val *vpp_stn.Rule) linuxclient.DataResyncDSL {
	defer func() { panic(deprecated) }()
	return dsl
}

// NAT44Global adds global NAT44 configuration to the RESYNC request.
func (dsl *DataResyncDSL) NAT44Global(val *vpp_nat.Nat44Global) linuxclient.DataResyncDSL {
	dsl.req.NatGlobal = val
	return dsl
}

// DNAT44 adds DNAT44 configuration to the RESYNC request
func (dsl *DataResyncDSL) DNAT44(val *vpp_nat.DNat44) linuxclient.DataResyncDSL {
	dsl.req.DNATs = append(dsl.req.DNATs, val)
	return dsl
}

// IPSecSA adds request to RESYNC a new Security Association
func (dsl *DataResyncDSL) IPSecSA(val *vpp_ipsec.SecurityAssociation) linuxclient.DataResyncDSL {
	dsl.req.SAs = append(dsl.req.SAs, val)
	return dsl
}

// IPSecSPD adds request to RESYNC a new Security Policy Database
func (dsl *DataResyncDSL) IPSecSPD(val *vpp_ipsec.SecurityPolicyDatabase) linuxclient.DataResyncDSL {
	dsl.req.SPDs = append(dsl.req.SPDs, val)
	return dsl
}

// PuntIPRedirect adds request to RESYNC a rule used to punt L3 traffic via interface.
func (dsl *DataResyncDSL) PuntIPRedirect(val *vpp_punt.IPRedirect) linuxclient.DataResyncDSL {
	dsl.req.IPRedirectPunts = append(dsl.req.IPRedirectPunts, val)
	return dsl
}

// PuntToHost adds request to RESYNC a rule used to punt L4 traffic to a host.
func (dsl *DataResyncDSL) PuntToHost(val *vpp_punt.ToHost) linuxclient.DataResyncDSL {
	dsl.req.ToHostPunts = append(dsl.req.ToHostPunts, val)
	return dsl
}

// Send propagates the request to the plugins. It deletes obsolete keys if listKeys() function is not null.
// The listkeys() function is used to list all current keys.
func (dsl *DataResyncDSL) Send() vppclient.Reply {
	var wasErr error

	ctx := context.Background()

	if _, err := dsl.client.Resync(ctx, dsl.req); err != nil {
		wasErr = err
	}

	return &Reply{err: wasErr}
}
