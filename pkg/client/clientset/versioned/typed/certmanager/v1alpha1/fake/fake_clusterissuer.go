/*
Copyright 2019 The Jetstack cert-manager contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/leki75/cert-manager/pkg/apis/certmanager/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterIssuers implements ClusterIssuerInterface
type FakeClusterIssuers struct {
	Fake *FakeCertmanagerV1alpha1
}

var clusterissuersResource = schema.GroupVersionResource{Group: "certmanager.k8s.io", Version: "v1alpha1", Resource: "clusterissuers"}

var clusterissuersKind = schema.GroupVersionKind{Group: "certmanager.k8s.io", Version: "v1alpha1", Kind: "ClusterIssuer"}

// Get takes name of the clusterIssuer, and returns the corresponding clusterIssuer object, and an error if there is any.
func (c *FakeClusterIssuers) Get(name string, options v1.GetOptions) (result *v1alpha1.ClusterIssuer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterissuersResource, name), &v1alpha1.ClusterIssuer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterIssuer), err
}

// List takes label and field selectors, and returns the list of ClusterIssuers that match those selectors.
func (c *FakeClusterIssuers) List(opts v1.ListOptions) (result *v1alpha1.ClusterIssuerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterissuersResource, clusterissuersKind, opts), &v1alpha1.ClusterIssuerList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ClusterIssuerList{ListMeta: obj.(*v1alpha1.ClusterIssuerList).ListMeta}
	for _, item := range obj.(*v1alpha1.ClusterIssuerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterIssuers.
func (c *FakeClusterIssuers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterissuersResource, opts))
}

// Create takes the representation of a clusterIssuer and creates it.  Returns the server's representation of the clusterIssuer, and an error, if there is any.
func (c *FakeClusterIssuers) Create(clusterIssuer *v1alpha1.ClusterIssuer) (result *v1alpha1.ClusterIssuer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterissuersResource, clusterIssuer), &v1alpha1.ClusterIssuer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterIssuer), err
}

// Update takes the representation of a clusterIssuer and updates it. Returns the server's representation of the clusterIssuer, and an error, if there is any.
func (c *FakeClusterIssuers) Update(clusterIssuer *v1alpha1.ClusterIssuer) (result *v1alpha1.ClusterIssuer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterissuersResource, clusterIssuer), &v1alpha1.ClusterIssuer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterIssuer), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeClusterIssuers) UpdateStatus(clusterIssuer *v1alpha1.ClusterIssuer) (*v1alpha1.ClusterIssuer, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(clusterissuersResource, "status", clusterIssuer), &v1alpha1.ClusterIssuer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterIssuer), err
}

// Delete takes name of the clusterIssuer and deletes it. Returns an error if one occurs.
func (c *FakeClusterIssuers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(clusterissuersResource, name), &v1alpha1.ClusterIssuer{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterIssuers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterissuersResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ClusterIssuerList{})
	return err
}

// Patch applies the patch and returns the patched clusterIssuer.
func (c *FakeClusterIssuers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterIssuer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterissuersResource, name, pt, data, subresources...), &v1alpha1.ClusterIssuer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterIssuer), err
}
