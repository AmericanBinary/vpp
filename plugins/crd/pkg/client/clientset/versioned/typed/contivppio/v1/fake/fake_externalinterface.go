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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	contivppiov1 "github.com/americanbinary/vpp/plugins/crd/pkg/apis/contivppio/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeExternalInterfaces implements ExternalInterfaceInterface
type FakeExternalInterfaces struct {
	Fake *FakeContivppV1
	ns   string
}

var externalinterfacesResource = schema.GroupVersionResource{Group: "contivpp.io", Version: "v1", Resource: "externalinterfaces"}

var externalinterfacesKind = schema.GroupVersionKind{Group: "contivpp.io", Version: "v1", Kind: "ExternalInterface"}

// Get takes name of the externalInterface, and returns the corresponding externalInterface object, and an error if there is any.
func (c *FakeExternalInterfaces) Get(name string, options v1.GetOptions) (result *contivppiov1.ExternalInterface, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(externalinterfacesResource, c.ns, name), &contivppiov1.ExternalInterface{})

	if obj == nil {
		return nil, err
	}
	return obj.(*contivppiov1.ExternalInterface), err
}

// List takes label and field selectors, and returns the list of ExternalInterfaces that match those selectors.
func (c *FakeExternalInterfaces) List(opts v1.ListOptions) (result *contivppiov1.ExternalInterfaceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(externalinterfacesResource, externalinterfacesKind, c.ns, opts), &contivppiov1.ExternalInterfaceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &contivppiov1.ExternalInterfaceList{ListMeta: obj.(*contivppiov1.ExternalInterfaceList).ListMeta}
	for _, item := range obj.(*contivppiov1.ExternalInterfaceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested externalInterfaces.
func (c *FakeExternalInterfaces) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(externalinterfacesResource, c.ns, opts))

}

// Create takes the representation of a externalInterface and creates it.  Returns the server's representation of the externalInterface, and an error, if there is any.
func (c *FakeExternalInterfaces) Create(externalInterface *contivppiov1.ExternalInterface) (result *contivppiov1.ExternalInterface, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(externalinterfacesResource, c.ns, externalInterface), &contivppiov1.ExternalInterface{})

	if obj == nil {
		return nil, err
	}
	return obj.(*contivppiov1.ExternalInterface), err
}

// Update takes the representation of a externalInterface and updates it. Returns the server's representation of the externalInterface, and an error, if there is any.
func (c *FakeExternalInterfaces) Update(externalInterface *contivppiov1.ExternalInterface) (result *contivppiov1.ExternalInterface, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(externalinterfacesResource, c.ns, externalInterface), &contivppiov1.ExternalInterface{})

	if obj == nil {
		return nil, err
	}
	return obj.(*contivppiov1.ExternalInterface), err
}

// Delete takes name of the externalInterface and deletes it. Returns an error if one occurs.
func (c *FakeExternalInterfaces) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(externalinterfacesResource, c.ns, name), &contivppiov1.ExternalInterface{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeExternalInterfaces) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(externalinterfacesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &contivppiov1.ExternalInterfaceList{})
	return err
}

// Patch applies the patch and returns the patched externalInterface.
func (c *FakeExternalInterfaces) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *contivppiov1.ExternalInterface, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(externalinterfacesResource, c.ns, name, pt, data, subresources...), &contivppiov1.ExternalInterface{})

	if obj == nil {
		return nil, err
	}
	return obj.(*contivppiov1.ExternalInterface), err
}
