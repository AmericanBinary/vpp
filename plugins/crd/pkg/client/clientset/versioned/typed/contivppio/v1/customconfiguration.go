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

package v1

import (
	"time"

	v1 "github.com/americanbinary/vpp/plugins/crd/pkg/apis/contivppio/v1"
	scheme "github.com/americanbinary/vpp/plugins/crd/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CustomConfigurationsGetter has a method to return a CustomConfigurationInterface.
// A group's client should implement this interface.
type CustomConfigurationsGetter interface {
	CustomConfigurations(namespace string) CustomConfigurationInterface
}

// CustomConfigurationInterface has methods to work with CustomConfiguration resources.
type CustomConfigurationInterface interface {
	Create(*v1.CustomConfiguration) (*v1.CustomConfiguration, error)
	Update(*v1.CustomConfiguration) (*v1.CustomConfiguration, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.CustomConfiguration, error)
	List(opts metav1.ListOptions) (*v1.CustomConfigurationList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CustomConfiguration, err error)
	CustomConfigurationExpansion
}

// customConfigurations implements CustomConfigurationInterface
type customConfigurations struct {
	client rest.Interface
	ns     string
}

// newCustomConfigurations returns a CustomConfigurations
func newCustomConfigurations(c *ContivppV1Client, namespace string) *customConfigurations {
	return &customConfigurations{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the customConfiguration, and returns the corresponding customConfiguration object, and an error if there is any.
func (c *customConfigurations) Get(name string, options metav1.GetOptions) (result *v1.CustomConfiguration, err error) {
	result = &v1.CustomConfiguration{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("customconfigurations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CustomConfigurations that match those selectors.
func (c *customConfigurations) List(opts metav1.ListOptions) (result *v1.CustomConfigurationList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.CustomConfigurationList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("customconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested customConfigurations.
func (c *customConfigurations) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("customconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a customConfiguration and creates it.  Returns the server's representation of the customConfiguration, and an error, if there is any.
func (c *customConfigurations) Create(customConfiguration *v1.CustomConfiguration) (result *v1.CustomConfiguration, err error) {
	result = &v1.CustomConfiguration{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("customconfigurations").
		Body(customConfiguration).
		Do().
		Into(result)
	return
}

// Update takes the representation of a customConfiguration and updates it. Returns the server's representation of the customConfiguration, and an error, if there is any.
func (c *customConfigurations) Update(customConfiguration *v1.CustomConfiguration) (result *v1.CustomConfiguration, err error) {
	result = &v1.CustomConfiguration{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("customconfigurations").
		Name(customConfiguration.Name).
		Body(customConfiguration).
		Do().
		Into(result)
	return
}

// Delete takes name of the customConfiguration and deletes it. Returns an error if one occurs.
func (c *customConfigurations) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("customconfigurations").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *customConfigurations) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("customconfigurations").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched customConfiguration.
func (c *customConfigurations) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CustomConfiguration, err error) {
	result = &v1.CustomConfiguration{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("customconfigurations").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
