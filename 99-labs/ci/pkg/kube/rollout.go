package kube

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// fatalWaitingReasons are container states that will not resolve by waiting: the pod spec
// or the image is wrong. Reporting them straight away turns a four minute timeout into a
// few seconds and an explanation.
var fatalWaitingReasons = map[string]bool{
	"CrashLoopBackOff":           true,
	"CreateContainerConfigError": true,
	"CreateContainerError":       true,
	"ErrImageNeverPull":          true,
	"ImagePullBackOff":           true,
	"InvalidImageName":           true,
	"RunContainerError":          true,
}

// fatalThreshold is how many consecutive observations of a fatal state we require, so that
// a container that merely restarts once on startup is not reported as broken.
const fatalThreshold = 3

// WaitRollout blocks until a Deployment has fully rolled out: the controller has observed
// the current spec, every replica has been replaced by an updated one, no old replica is
// left, and all of them are available.
//
// This is what "kubectl rollout status" checks, and the strictness matters: a plain "are
// enough pods ready" test passes immediately after a restart is requested, while the old
// pod is still running and the new one does not exist yet, and the caller then talks to
// the very pod it was trying to replace.
func (c *Client) WaitRollout(ctx context.Context, name string, timeout time.Duration) error {
	fatal := 0

	err := poll(ctx, timeout, func() (bool, string, error) {
		d, err := c.Clientset.AppsV1().Deployments(c.Namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, "not found yet", ignoreNotFound(err)
		}

		if selector := selectorString(d.Spec.Selector); selector != "" {
			if reason := c.fatalPodReason(ctx, selector); reason != "" {
				fatal++
				if fatal >= fatalThreshold {
					return false, "", fmt.Errorf("the pods cannot start: %s", reason)
				}
			} else {
				fatal = 0
			}
		}

		if d.Status.ObservedGeneration < d.Generation {
			return false, "the Deployment controller has not observed the update yet", nil
		}

		want := int32(1)
		if d.Spec.Replicas != nil {
			want = *d.Spec.Replicas
		}

		switch {
		case d.Status.UpdatedReplicas < want:
			return false, fmt.Sprintf("%d/%d replicas updated",
				d.Status.UpdatedReplicas, want), nil
		case d.Status.Replicas > d.Status.UpdatedReplicas:
			return false, fmt.Sprintf("%d replicas of the previous version still running",
				d.Status.Replicas-d.Status.UpdatedReplicas), nil
		case d.Status.AvailableReplicas < want:
			return false, fmt.Sprintf("%d/%d replicas available",
				d.Status.AvailableReplicas, want), nil
		}

		return true, "", nil
	})

	if err != nil {
		return fmt.Errorf("Deployment %q did not roll out: %w", name, err)
	}

	return nil
}

// fatalPodReason reports why the pods matching a selector cannot start, or an empty string
// when nothing is obviously wrong.
func (c *Client) fatalPodReason(ctx context.Context, selector string) string {
	pods, err := c.Clientset.CoreV1().Pods(c.Namespace).List(ctx,
		metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return ""
	}

	for i := range pods.Items {
		p := &pods.Items[i]
		if p.DeletionTimestamp != nil {
			continue
		}

		for _, cs := range p.Status.ContainerStatuses {
			w := cs.State.Waiting
			if w == nil || !fatalWaitingReasons[w.Reason] {
				continue
			}

			msg := fmt.Sprintf("pod %s, container %s: %s", p.Name, cs.Name, w.Reason)
			if w.Message != "" {
				msg += ": " + w.Message
			}

			return msg
		}
	}

	return ""
}

// selectorString renders a label selector in the form the list API expects.
func selectorString(s *metav1.LabelSelector) string {
	if s == nil {
		return ""
	}

	sel, err := metav1.LabelSelectorAsSelector(s)
	if err != nil {
		return ""
	}

	return sel.String()
}

// PodsForDeployment returns the pods belonging to a Deployment.
func (c *Client) PodsForDeployment(ctx context.Context, name string) ([]corev1.Pod, error) {
	d, err := c.Deployment(ctx, name)
	if err != nil {
		return nil, err
	}

	sel, err := metav1.LabelSelectorAsMap(d.Spec.Selector)
	if err != nil {
		return nil, fmt.Errorf("Deployment %q has an unusable selector: %w", name, err)
	}

	pods, err := c.Clientset.CoreV1().Pods(c.Namespace).List(ctx,
		metav1.ListOptions{LabelSelector: labels.SelectorFromSet(sel).String()})
	if err != nil {
		return nil, fmt.Errorf("list pods of Deployment %q: %w", name, err)
	}

	return pods.Items, nil
}
