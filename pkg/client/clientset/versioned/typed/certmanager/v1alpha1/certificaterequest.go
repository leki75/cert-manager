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

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/leki75/cert-manager/pkg/apis/certmanager/v1alpha1"
	scheme "github.com/leki75/cert-manager/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CertificateRequestsGetter has a method to return a CertificateRequestInterface.
// A group's client should implement this interface.
type CertificateRequestsGetter interface {
	CertificateRequests(namespace string) CertificateRequestInterface
}

// CertificateRequestInterface has methods to work with CertificateRequest resources.
type CertificateRequestInterface interface {
	Create(*v1alpha1.CertificateRequest) (*v1alpha1.CertificateRequest, error)
	Update(*v1alpha1.CertificateRequest) (*v1alpha1.CertificateRequest, error)
	UpdateStatus(*v1alpha1.CertificateRequest) (*v1alpha1.CertificateRequest, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.CertificateRequest, error)
	List(opts v1.ListOptions) (*v1alpha1.CertificateRequestList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CertificateRequest, err error)
	CertificateRequestExpansion
}

// certificateRequests implements CertificateRequestInterface
type certificateRequests struct {
	client rest.Interface
	ns     string
}

// newCertificateRequests returns a CertificateRequests
func newCertificateRequests(c *CertmanagerV1alpha1Client, namespace string) *certificateRequests {
	return &certificateRequests{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the certificateRequest, and returns the corresponding certificateRequest object, and an error if there is any.
func (c *certificateRequests) Get(name string, options v1.GetOptions) (result *v1alpha1.CertificateRequest, err error) {
	result = &v1alpha1.CertificateRequest{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("certificaterequests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CertificateRequests that match those selectors.
func (c *certificateRequests) List(opts v1.ListOptions) (result *v1alpha1.CertificateRequestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.CertificateRequestList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("certificaterequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested certificateRequests.
func (c *certificateRequests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("certificaterequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a certificateRequest and creates it.  Returns the server's representation of the certificateRequest, and an error, if there is any.
func (c *certificateRequests) Create(certificateRequest *v1alpha1.CertificateRequest) (result *v1alpha1.CertificateRequest, err error) {
	result = &v1alpha1.CertificateRequest{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("certificaterequests").
		Body(certificateRequest).
		Do().
		Into(result)
	return
}

// Update takes the representation of a certificateRequest and updates it. Returns the server's representation of the certificateRequest, and an error, if there is any.
func (c *certificateRequests) Update(certificateRequest *v1alpha1.CertificateRequest) (result *v1alpha1.CertificateRequest, err error) {
	result = &v1alpha1.CertificateRequest{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("certificaterequests").
		Name(certificateRequest.Name).
		Body(certificateRequest).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *certificateRequests) UpdateStatus(certificateRequest *v1alpha1.CertificateRequest) (result *v1alpha1.CertificateRequest, err error) {
	result = &v1alpha1.CertificateRequest{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("certificaterequests").
		Name(certificateRequest.Name).
		SubResource("status").
		Body(certificateRequest).
		Do().
		Into(result)
	return
}

// Delete takes name of the certificateRequest and deletes it. Returns an error if one occurs.
func (c *certificateRequests) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("certificaterequests").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *certificateRequests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("certificaterequests").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched certificateRequest.
func (c *certificateRequests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CertificateRequest, err error) {
	result = &v1alpha1.CertificateRequest{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("certificaterequests").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
