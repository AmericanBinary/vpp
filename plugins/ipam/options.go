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

package ipam

import (
	"go.ligato.io/cn-infra/v2/logging"

	"github.com/americanbinary/vpp/plugins/contivconf"
	"github.com/americanbinary/vpp/plugins/nodesync"
	"go.ligato.io/cn-infra/v2/rpc/rest"
	"go.ligato.io/cn-infra/v2/servicelabel"
)

// DefaultPlugin is a default instance of IPAM plugin.
var DefaultPlugin = *NewPlugin()

// NewPlugin creates a new Plugin with the provides Options
func NewPlugin(opts ...Option) *IPAM {
	p := &IPAM{}

	p.PluginName = "ipam"
	p.NodeSync = &nodesync.DefaultPlugin
	p.ContivConf = &contivconf.DefaultPlugin
	p.ServiceLabel = &servicelabel.DefaultPlugin
	p.HTTPHandlers = &rest.DefaultPlugin

	for _, o := range opts {
		o(p)
	}

	if p.Deps.Log == nil {
		p.Deps.Log = logging.ForPlugin(p.String())
	}

	return p
}

// Option is a function that acts on a Plugin to inject Dependencies or configuration
type Option func(*IPAM)

// UseDeps returns Option that can inject custom dependencies.
func UseDeps(cb func(*Deps)) Option {
	return func(p *IPAM) {
		cb(&p.Deps)
	}
}
