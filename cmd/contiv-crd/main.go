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

package main

import (
	"github.com/americanbinary/vpp/plugins/crd"
	"github.com/americanbinary/vpp/plugins/ksr"
	"go.ligato.io/cn-infra/v2/agent"
	"go.ligato.io/cn-infra/v2/datasync/kvdbsync"
	"go.ligato.io/cn-infra/v2/datasync/resync"
	"go.ligato.io/cn-infra/v2/db/keyval/etcd"
	"go.ligato.io/cn-infra/v2/health/probe"
	"go.ligato.io/cn-infra/v2/logging/logrus"
	"go.ligato.io/cn-infra/v2/servicelabel"

	// load all VPP-agent models for CustomConfiguration CRD handler to use
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/linux"
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/linux/iptables"
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/vpp"
	_ "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/srv6"
)

// ContivCRD is a custom resource to provide Contiv-VPP telemetry information.
type ContivCRD struct {
	HealthProbe *probe.Plugin
	CRD         *crd.Plugin
}

func (c *ContivCRD) String() string {
	return "CRD"
}

// Init is called at startup phase. Method added in order to implement Plugin interface.
func (c *ContivCRD) Init() error {
	return nil
}

// AfterInit triggers the first resync.
func (c *ContivCRD) AfterInit() error {
	resync.DefaultPlugin.DoResync()
	return nil
}

// Close is called at cleanup phase. Method added in order to implement Plugin interface.
func (c *ContivCRD) Close() error {
	return nil
}

func main() {
	ksrServicelabel := servicelabel.NewPlugin(servicelabel.UseLabel(ksr.MicroserviceLabel))
	ksrServicelabel.SetName("ksrServiceLabel")

	ksrDataSync := kvdbsync.NewPlugin(kvdbsync.UseDeps(func(deps *kvdbsync.Deps) {
		deps.KvPlugin = &etcd.DefaultPlugin
		deps.ServiceLabel = ksrServicelabel
		deps.SetName("ksrDataSync")
	}))

	// disable status check for etcd
	etcd.DefaultPlugin.StatusCheck = nil

	crd.DefaultPlugin.Watcher = ksrDataSync
	crd.DefaultPlugin.Etcd = &etcd.DefaultPlugin
	crd.DefaultPlugin.ServiceLabel = ksrServicelabel

	probe.DefaultPlugin.NonFatalPlugins = []string{"etcd"}

	ContivCRD := &ContivCRD{
		HealthProbe: &probe.DefaultPlugin,
		CRD:         &crd.DefaultPlugin,
	}

	a := agent.NewAgent(agent.AllPlugins(ContivCRD))
	if err := a.Run(); err != nil {
		logrus.DefaultLogger().Fatal(err)
	}
}
