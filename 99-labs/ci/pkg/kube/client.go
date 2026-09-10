// Package kube is a thin Kubernetes layer for the lab tests: it applies the manifests a
// student wrote, waits until the workloads are actually serving, inspects what they
// declared, and reports useful diagnostics when something does not come up.
package kube

import (
	"context"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

// Client bundles the Kubernetes clients the lab tests need.
type Client struct {
	Clientset *kubernetes.Clientset
	Dynamic   dynamic.Interface
	Mapper    meta.RESTMapper
	Namespace string
}

// NewClient builds a client from the kubeconfig named in the KUBECONFIG environment
// variable, falling back to the usual default location.
func NewClient() (*Client, error) {
	path := os.Getenv("KUBECONFIG")
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("no KUBECONFIG set and no home directory: %w", err)
		}
		path = home + "/.kube/config"
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, fmt.Errorf("load kubeconfig %q: %w", path, err)
	}

	// The lab clusters are small and the tests are chatty.
	cfg.QPS, cfg.Burst = 50, 100

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build clientset: %w", err)
	}

	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build dynamic client: %w", err)
	}

	disco, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("build discovery client: %w", err)
	}

	groups, err := restmapper.GetAPIGroupResources(disco)
	if err != nil {
		return nil, fmt.Errorf("discover API groups: %w", err)
	}

	return &Client{
		Clientset: cs,
		Dynamic:   dyn,
		Mapper:    restmapper.NewDiscoveryRESTMapper(groups),
		Namespace: "default",
	}, nil
}

// resourceFor returns the dynamic client for an object, honouring whether the object's
// kind is namespaced.
func (c *Client) resourceFor(obj *Object) (dynamic.ResourceInterface, error) {
	gvk := obj.GroupVersionKind()
	mapping, err := c.Mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, fmt.Errorf("no REST mapping for %s: %w", gvk, err)
	}

	if mapping.Scope.Name() == meta.RESTScopeNameRoot {
		return c.Dynamic.Resource(mapping.Resource), nil
	}

	ns := obj.GetNamespace()
	if ns == "" {
		ns = c.Namespace
	}

	return c.Dynamic.Resource(mapping.Resource).Namespace(ns), nil
}

// Ping checks that the API server answers, so that a broken cluster fails fast and with a
// clear message rather than as a pile of timeouts.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("cannot reach the Kubernetes API server: %w", err)
	}
	return nil
}
