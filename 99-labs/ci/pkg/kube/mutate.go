package kube

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// SetContainerEnv sets an environment variable on a container in a workload manifest,
// before it is applied. The CI uses it to switch on facilities that a student build must
// not have, such as fault injection in the key-value store.
func SetContainerEnv(obj *Object, container, name, value string) error {
	path := []string{"spec", "template", "spec", "containers"}

	ctrs, found, err := unstructured.NestedSlice(obj.Object, path...)
	if err != nil {
		return fmt.Errorf("read containers of %s/%s: %w", obj.GetKind(), obj.GetName(), err)
	}
	if !found {
		return fmt.Errorf("%s/%s has no pod template", obj.GetKind(), obj.GetName())
	}

	patched := false
	for i, raw := range ctrs {
		c, ok := raw.(map[string]any)
		if !ok || c["name"] != container {
			continue
		}

		env, _ := c["env"].([]any)

		replaced := false
		for j, rawVar := range env {
			v, ok := rawVar.(map[string]any)
			if !ok || v["name"] != name {
				continue
			}
			env[j] = map[string]any{"name": name, "value": value}
			replaced = true
		}
		if !replaced {
			env = append(env, map[string]any{"name": name, "value": value})
		}

		c["env"] = env
		ctrs[i] = c
		patched = true
	}

	if !patched {
		return fmt.Errorf("%s/%s has no container named %q",
			obj.GetKind(), obj.GetName(), container)
	}

	return unstructured.SetNestedSlice(obj.Object, ctrs, path...)
}

// DeletePVCs removes the persistent volume claims matching a label selector. A
// StatefulSet's volumes outlive the StatefulSet, so a teardown that skipped them would
// leak one lab's data into the next.
func (c *Client) DeletePVCs(ctx context.Context, selector string) error {
	list, err := c.Clientset.CoreV1().PersistentVolumeClaims(c.Namespace).
		List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return fmt.Errorf("list PVCs %q: %w", selector, err)
	}

	for _, pvc := range list.Items {
		err := c.Clientset.CoreV1().PersistentVolumeClaims(c.Namespace).
			Delete(ctx, pvc.Name, metav1.DeleteOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			return fmt.Errorf("delete PVC %q: %w", pvc.Name, err)
		}
	}

	// Deletion is asynchronous, and a StatefulSet that is recreated while its claim is
	// still going away fails to start its pod at all.
	return poll(ctx, 2*time.Minute, func() (bool, string, error) {
		left, err := c.Clientset.CoreV1().PersistentVolumeClaims(c.Namespace).
			List(ctx, metav1.ListOptions{LabelSelector: selector})
		if err != nil {
			return false, "", err
		}
		if len(left.Items) == 0 {
			return true, "", nil
		}

		return false, fmt.Sprintf("%d volume claims still being deleted", len(left.Items)), nil
	})
}

// localImagePrefix marks an image that was built into the cluster's own image store rather
// than pushed to a registry, which is how the labs work.
const localImagePrefix = "localhost/"

// NormalizeLocalImages makes sure locally built images are not pulled from a registry, and
// reports what it changed.
//
// The labs are written for Minikube with the cri-o runtime, where the "localhost/" prefix
// names the node's local image storage. Under containerd, which is what the CI cluster
// runs, the same prefix names a registry host called "localhost", so the default pull
// policy for a ":latest" image sends the kubelet off to a registry that does not exist.
// Rewriting the policy compensates for that difference between the runtimes; it is not
// papering over a mistake in the manifest.
func NormalizeLocalImages(obj *Object) ([]string, error) {
	changed := []string{}

	for _, field := range []string{"containers", "initContainers"} {
		path := []string{"spec", "template", "spec", field}

		ctrs, found, err := unstructured.NestedSlice(obj.Object, path...)
		if err != nil || !found {
			continue
		}

		patched := false
		for i, raw := range ctrs {
			c, ok := raw.(map[string]any)
			if !ok {
				continue
			}

			image, _ := c["image"].(string)
			if !strings.HasPrefix(image, localImagePrefix) {
				continue
			}

			policy, _ := c["imagePullPolicy"].(string)
			if policy == string(corev1.PullNever) || policy == string(corev1.PullIfNotPresent) {
				continue
			}

			c["imagePullPolicy"] = string(corev1.PullIfNotPresent)
			ctrs[i] = c
			patched = true

			changed = append(changed, fmt.Sprintf("%s/%s: container %v uses the local image %q, "+
				"pull policy set to IfNotPresent", obj.GetKind(), obj.GetName(), c["name"], image))
		}

		if patched {
			if err := unstructured.SetNestedSlice(obj.Object, ctrs, path...); err != nil {
				return nil, fmt.Errorf("%s/%s: %w", obj.GetKind(), obj.GetName(), err)
			}
		}
	}

	return changed, nil
}
