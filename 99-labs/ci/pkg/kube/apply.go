package kube

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// fieldManager identifies our writes to the API server for server-side apply.
const fieldManager = "labci"

// Object is a Kubernetes object read from a manifest.
type Object = unstructured.Unstructured

// Parse splits a multi-document YAML manifest into individual objects.
func Parse(manifest []byte) ([]*Object, error) {
	dec := yaml.NewYAMLOrJSONDecoder(strings.NewReader(string(manifest)), 4096)

	objs := []*Object{}
	for {
		obj := &Object{}
		err := dec.Decode(obj)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse manifest: %w", err)
		}

		// A trailing "---" or a comment-only document decodes to nothing.
		if len(obj.Object) == 0 {
			continue
		}

		if obj.GetKind() == "" {
			return nil, fmt.Errorf("parse manifest: object without a kind: %v", obj.Object)
		}

		objs = append(objs, obj)
	}

	return objs, nil
}

// ParseFile reads and splits a manifest file.
func ParseFile(path string) ([]*Object, error) {
	manifest, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}

	objs, err := Parse(manifest)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	if len(objs) == 0 {
		return nil, fmt.Errorf("%s: manifest contains no Kubernetes objects", path)
	}

	return objs, nil
}

// Apply server-side applies the objects, in the order given. Objects without a namespace
// land in the client's default namespace.
func (c *Client) Apply(ctx context.Context, objs []*Object) error {
	for _, obj := range objs {
		ri, err := c.resourceFor(obj)
		if err != nil {
			return err
		}

		if obj.GetNamespace() == "" && c.namespaced(obj) {
			obj.SetNamespace(c.Namespace)
		}

		if _, err := ri.Apply(ctx, obj.GetName(), obj, metav1.ApplyOptions{
			FieldManager: fieldManager,
			Force:        true,
		}); err != nil {
			return fmt.Errorf("apply %s/%s: %w", obj.GetKind(), obj.GetName(), err)
		}
	}

	return nil
}

// ApplyFile applies a manifest file and returns the objects it contained.
func (c *Client) ApplyFile(ctx context.Context, path string) ([]*Object, error) {
	objs, err := ParseFile(path)
	if err != nil {
		return nil, err
	}

	return objs, c.Apply(ctx, objs)
}

// ApplyYAML applies an inline manifest and returns the objects it contained.
func (c *Client) ApplyYAML(ctx context.Context, manifest string) ([]*Object, error) {
	objs, err := Parse([]byte(manifest))
	if err != nil {
		return nil, err
	}

	return objs, c.Apply(ctx, objs)
}

// Delete removes the objects, ignoring the ones that are already gone.
func (c *Client) Delete(ctx context.Context, objs []*Object) error {
	// Delete in reverse order so that workloads go before the config they reference.
	for i := len(objs) - 1; i >= 0; i-- {
		obj := objs[i]

		ri, err := c.resourceFor(obj)
		if err != nil {
			// An unknown kind cannot have been applied either.
			continue
		}

		err = ri.Delete(ctx, obj.GetName(), metav1.DeleteOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("delete %s/%s: %w", obj.GetKind(), obj.GetName(), err)
		}
	}

	return nil
}

// namespaced reports whether the object's kind lives in a namespace.
func (c *Client) namespaced(obj *Object) bool {
	gvk := obj.GroupVersionKind()
	mapping, err := c.Mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return true
	}
	return mapping.Scope.Name() != "root"
}
