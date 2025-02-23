// Copyright (c) 2018 Cisco and/or its affiliates.
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

package cache

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.ligato.io/cn-infra/v2/datasync"

	nodemodel "github.com/americanbinary/vpp/plugins/ksr/model/node"
	podmodel "github.com/americanbinary/vpp/plugins/ksr/model/pod"
	vppnodemodel "github.com/americanbinary/vpp/plugins/nodesync/vppnode"

	"github.com/americanbinary/vpp/plugins/crd/api"
)

// Resync sends the resync event passed as an argument to the ctc telemetryCache
// thread, where it processed in the function below (resync)
func (ctc *ContivTelemetryCache) Resync(resyncEv datasync.ResyncEvent) error {
	ctc.dsUpdateChannel <- resyncEv
	return nil
}

// resync is an internal function that processes a data sync re sync event
// associated with K8s State data. The function executes in context of the
// ctc telemetryCache thread. The cache content is fully replaced with the received
// data.
func (ctc *ContivTelemetryCache) resync(resyncEv datasync.ResyncEvent) error {
	err := error(nil)
	ctc.Synced = true

	// Delete all previous data from cache, we are starting from scratch again
	ctc.ReinitializeCache()

	for resyncKey, resyncData := range resyncEv.GetValues() {
		for {
			evData, stop := resyncData.GetNext()
			if stop {
				break
			}

			key := evData.GetKey()
			switch resyncKey {
			case vppnodemodel.KeyPrefix:
				err = ctc.parseAndCacheVppNodeData(key, evData)

			case podmodel.KeyPrefix():
				err = ctc.parseAndCachePodData(key, evData)

			case nodemodel.KeyPrefix():
				err = ctc.parseAndCacheNodeData(key, evData)

			default:
				err = fmt.Errorf("unknown RESYNC Key %s, key %s", resyncKey, key)
			}

			if err != nil {
				ctc.Report.LogErrAndAppendToNodeReport(api.GlobalMsg, err.Error())
				ctc.Synced = false
			}
		}
	}

	if ctc.Synced == false {
		retErr := fmt.Errorf("datasync error, cache may be out of sync")
		ctc.Report.AppendToNodeReport(api.GlobalMsg, retErr.Error())
		return retErr

	}
	return nil
}

func (ctc *ContivTelemetryCache) parseAndCacheVppNodeData(key string, evData datasync.KeyVal) error {
	pattern := fmt.Sprintf("%s[0-9]*$", vppnodemodel.KeyPrefix)
	matched, err := regexp.Match(pattern, []byte(key))
	if !matched || err != nil {
		return fmt.Errorf("invalid key %s", key)
	}

	vppNodeValue := &vppnodemodel.VppNode{}
	err = evData.GetValue(vppNodeValue)
	if err != nil {
		return fmt.Errorf("could not parse node info data for key %s, error %s", key, err)
	}

	ctc.Log.Infof("parseAndCacheVppNodeData: node %s, id: %d", vppNodeValue.Name, vppNodeValue.Id)

	id, _ := strconv.Atoi(strings.Split(key, "/")[1])
	if vppNodeValue.Id != uint32(id) {
		return fmt.Errorf("invalid key '%s' or node id '%d'", key, vppNodeValue.Id)
	}

	if vppNodeValue.Id == 0 || vppNodeValue.Name == "" || len(vppNodeValue.IpAddresses) == 0 {
		return fmt.Errorf("invalid nodeInfo data: '%+v'", vppNodeValue)
	}

	err = ctc.VppCache.CreateNode(vppNodeValue.Id, vppNodeValue.Name,
		vppNodeValue.IpAddresses[0])
	if err != nil {
		return fmt.Errorf("failed to add vpp node, error: '%s'", err)
	}

	return nil
}

func (ctc *ContivTelemetryCache) parseAndCachePodData(key string, evData datasync.KeyVal) error {
	pod, namespace, err := podmodel.ParsePodFromKey(key)
	if err != nil {
		return fmt.Errorf("invalid key %s", key)
	}

	podValue := &podmodel.Pod{}
	err = evData.GetValue(podValue)
	if err != nil {
		return fmt.Errorf("could not parse node info data for key %s, error %s", key, err)
	}

	ctc.Log.Infof("parseAndCachePodData: pod %s, namespace %s, value %+v",
		pod, namespace, podValue)
	ctc.K8sCache.CreatePod(podValue.Name, podValue.Namespace, podValue.Label, podValue.IpAddress,
		podValue.HostIpAddress, podValue.Container)

	return nil
}

func (ctc *ContivTelemetryCache) parseAndCacheNodeData(key string, evData datasync.KeyVal) error {
	node, err := nodemodel.ParseNodeFromKey(key)
	if err != nil {
		return fmt.Errorf("invalid key %s", key)
	}

	nodeValue := &nodemodel.Node{}
	err = evData.GetValue(nodeValue)
	if err != nil {
		return fmt.Errorf("could not parse node info data for key %s, error %s", key, err)
	}

	ctc.Log.Infof("parseAndCacheNodeData: node %s, value %+v", node, nodeValue)
	ctc.K8sCache.CreateK8sNode(nodeValue.Name, nodeValue.Pod_CIDR, nodeValue.Provider_ID,
		nodeValue.Addresses, nodeValue.NodeInfo)
	return nil
}
