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
	"github.com/americanbinary/vpp/plugins/ksr"
	"go.ligato.io/cn-infra/v2/agent"
	"go.ligato.io/cn-infra/v2/datasync/kvdbsync"
	"go.ligato.io/cn-infra/v2/datasync/resync"
	"go.ligato.io/cn-infra/v2/db/keyval/etcd"
	"go.ligato.io/cn-infra/v2/health/probe"
	"go.ligato.io/cn-infra/v2/logging/logrus"
	"go.ligato.io/cn-infra/v2/servicelabel"
)

// ContivKSR reflects kubernetes data into contiv/vpp's etcd.
type ContivKSR struct {
	ServiceLabel servicelabel.ReaderAPI
	HealthProbe  *probe.Plugin
	DataSyncETCD *kvdbsync.Plugin
	KSR          *ksr.Plugin
}

func (c *ContivKSR) String() string {
	return "KSR"
}

// Init is called at startup phase. Method added in order to implement Plugin interface.
func (c *ContivKSR) Init() error {
	return nil
}

// Close is called at cleanup phase. Method added in order to implement Plugin interface.
func (c *ContivKSR) Close() error {
	return nil
}

func main() {

	servicelabel.DefaultPlugin.MicroserviceLabel = ksr.MicroserviceLabel

	etcdDataSync := kvdbsync.NewPlugin(kvdbsync.UseDeps(func(deps *kvdbsync.Deps) {
		deps.KvPlugin = &etcd.DefaultPlugin
		deps.ResyncOrch = &resync.DefaultPlugin
	}))

	ksr.DefaultPlugin.Publish = etcdDataSync

	probe.DefaultPlugin.NonFatalPlugins = []string{"etcd"}

	contivKSR := &ContivKSR{
		ServiceLabel: &servicelabel.DefaultPlugin,
		HealthProbe:  &probe.DefaultPlugin,
		DataSyncETCD: etcdDataSync,
		KSR:          &ksr.DefaultPlugin,
	}

	a := agent.NewAgent(agent.AllPlugins(contivKSR))
	if err := a.Run(); err != nil {
		logrus.DefaultLogger().Fatal(err)
	}
}
