package kube

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// The helpers below read back what a student declared in their manifests. Several lab
// requirements (a liveness probe, a termination grace period, a ConfigMap wired into an
// environment variable) are not observable from the app's HTTP API, so the tests check the
// declaration and the resulting behaviour separately.

// Deployment fetches a Deployment.
func (c *Client) Deployment(ctx context.Context, name string) (*appsv1.Deployment, error) {
	d, err := c.Clientset.AppsV1().Deployments(c.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get Deployment %q: %w", name, err)
	}
	return d, nil
}

// Container returns the named container from a pod spec, or the only container when the
// name is empty.
func Container(spec *corev1.PodSpec, name string) (*corev1.Container, error) {
	if name == "" && len(spec.Containers) == 1 {
		return &spec.Containers[0], nil
	}

	for i := range spec.Containers {
		if spec.Containers[i].Name == name {
			return &spec.Containers[i], nil
		}
	}

	return nil, fmt.Errorf("no container %q in the pod spec", name)
}

// ConfigMapRef describes how a container sources an environment variable from a ConfigMap.
type ConfigMapRef struct {
	// EnvVar is the environment variable being set.
	EnvVar string
	// Name is the ConfigMap it reads from.
	Name string
	// Key is the entry within that ConfigMap.
	Key string
}

// ConfigMapRefFor reports how a container obtains an environment variable from a
// ConfigMap, or nil when the variable is not wired to one.
func ConfigMapRefFor(ctr *corev1.Container, envVar string) *ConfigMapRef {
	for _, e := range ctr.Env {
		if e.Name != envVar {
			continue
		}
		if e.ValueFrom == nil || e.ValueFrom.ConfigMapKeyRef == nil {
			return nil
		}

		return &ConfigMapRef{
			EnvVar: envVar,
			Name:   e.ValueFrom.ConfigMapKeyRef.Name,
			Key:    e.ValueFrom.ConfigMapKeyRef.Key,
		}
	}

	return nil
}

// Scale sets the replica count of a Deployment.
func (c *Client) Scale(ctx context.Context, name string, replicas int32) error {
	s, err := c.Clientset.AppsV1().Deployments(c.Namespace).GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get scale of Deployment %q: %w", name, err)
	}

	s.Spec.Replicas = replicas
	if _, err := c.Clientset.AppsV1().Deployments(c.Namespace).
		UpdateScale(ctx, name, s, metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("scale Deployment %q to %d: %w", name, replicas, err)
	}

	return nil
}

// ScaleStatefulSet sets the replica count of a StatefulSet. The lab tests use it to take
// the key-value store away from under the app.
func (c *Client) ScaleStatefulSet(ctx context.Context, name string, replicas int32) error {
	s, err := c.Clientset.AppsV1().StatefulSets(c.Namespace).GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get scale of StatefulSet %q: %w", name, err)
	}

	s.Spec.Replicas = replicas
	if _, err := c.Clientset.AppsV1().StatefulSets(c.Namespace).
		UpdateScale(ctx, name, s, metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("scale StatefulSet %q to %d: %w", name, replicas, err)
	}

	return nil
}

// restartPatch is the annotation kubectl uses to trigger a rollout restart.
const restartPatch = `{"spec":{"template":{"metadata":{"annotations":{"labci/restartedAt":%q}}}}}`

// RolloutRestart restarts a Deployment the way "kubectl rollout restart" does.
func (c *Client) RolloutRestart(ctx context.Context, name string) error {
	patch := fmt.Sprintf(restartPatch, time.Now().Format(time.RFC3339Nano))

	if _, err := c.Clientset.AppsV1().Deployments(c.Namespace).
		Patch(ctx, name, types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{}); err != nil {
		return fmt.Errorf("restart Deployment %q: %w", name, err)
	}

	return nil
}

// RolloutRestartStatefulSet restarts a StatefulSet.
func (c *Client) RolloutRestartStatefulSet(ctx context.Context, name string) error {
	patch := fmt.Sprintf(restartPatch, time.Now().Format(time.RFC3339Nano))

	if _, err := c.Clientset.AppsV1().StatefulSets(c.Namespace).
		Patch(ctx, name, types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{}); err != nil {
		return fmt.Errorf("restart StatefulSet %q: %w", name, err)
	}

	return nil
}

// SetEnv sets an environment variable on a container in a Deployment, replacing any
// existing definition of that variable including one sourced from a ConfigMap.
func (c *Client) SetEnv(ctx context.Context, deployment, container, name, value string) error {
	d, err := c.Deployment(ctx, deployment)
	if err != nil {
		return err
	}

	ctr, err := Container(&d.Spec.Template.Spec, container)
	if err != nil {
		return fmt.Errorf("Deployment %q: %w", deployment, err)
	}

	env := []any{}
	replaced := false
	for _, e := range ctr.Env {
		if e.Name == name {
			env = append(env, map[string]any{"name": name, "value": value})
			replaced = true
			continue
		}
		entry := map[string]any{"name": e.Name}
		if e.Value != "" {
			entry["value"] = e.Value
		}
		env = append(env, entry)
	}
	if !replaced {
		env = append(env, map[string]any{"name": name, "value": value})
	}

	// A strategic merge on "env" merges by key, so a variable that used to have a
	// valueFrom would keep it. Replace the whole list instead.
	patch := map[string]any{
		"spec": map[string]any{
			"template": map[string]any{
				"spec": map[string]any{
					"containers": []any{
						map[string]any{"name": ctr.Name, "env": env},
					},
				},
			},
		},
	}

	return c.patchDeployment(ctx, deployment, patch)
}

func (c *Client) patchDeployment(ctx context.Context, name string, patch map[string]any) error {
	body, err := marshal(patch)
	if err != nil {
		return err
	}

	if _, err := c.Clientset.AppsV1().Deployments(c.Namespace).
		Patch(ctx, name, types.StrategicMergePatchType, body, metav1.PatchOptions{}); err != nil {
		return fmt.Errorf("patch Deployment %q: %w", name, err)
	}

	return nil
}

// DeletePod removes a single pod, giving it the requested grace period to shut down.
func (c *Client) DeletePod(ctx context.Context, name string, grace int64) error {
	err := c.Clientset.CoreV1().Pods(c.Namespace).Delete(ctx, name, metav1.DeleteOptions{
		GracePeriodSeconds: &grace,
	})
	if err != nil {
		return fmt.Errorf("delete pod %q: %w", name, err)
	}

	return nil
}

// ConfigMapValue returns an entry of a ConfigMap.
func (c *Client) ConfigMapValue(ctx context.Context, name, key string) (string, bool, error) {
	cm, err := c.Clientset.CoreV1().ConfigMaps(c.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", false, fmt.Errorf("get ConfigMap %q: %w", name, err)
	}

	v, ok := cm.Data[key]
	return v, ok, nil
}
