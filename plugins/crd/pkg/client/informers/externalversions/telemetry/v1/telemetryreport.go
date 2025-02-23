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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	telemetryv1 "github.com/americanbinary/vpp/plugins/crd/pkg/apis/telemetry/v1"
	versioned "github.com/americanbinary/vpp/plugins/crd/pkg/client/clientset/versioned"
	internalinterfaces "github.com/americanbinary/vpp/plugins/crd/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/americanbinary/vpp/plugins/crd/pkg/client/listers/telemetry/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TelemetryReportInformer provides access to a shared informer and lister for
// TelemetryReports.
type TelemetryReportInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.TelemetryReportLister
}

type telemetryReportInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTelemetryReportInformer constructs a new informer for TelemetryReport type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTelemetryReportInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTelemetryReportInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTelemetryReportInformer constructs a new informer for TelemetryReport type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTelemetryReportInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TelemetryV1().TelemetryReports(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TelemetryV1().TelemetryReports(namespace).Watch(options)
			},
		},
		&telemetryv1.TelemetryReport{},
		resyncPeriod,
		indexers,
	)
}

func (f *telemetryReportInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTelemetryReportInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *telemetryReportInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&telemetryv1.TelemetryReport{}, f.defaultInformer)
}

func (f *telemetryReportInformer) Lister() v1.TelemetryReportLister {
	return v1.NewTelemetryReportLister(f.Informer().GetIndexer())
}
