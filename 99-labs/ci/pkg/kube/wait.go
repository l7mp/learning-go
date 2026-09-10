package kube

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// pollInterval is how often the wait helpers re-check the API server.
const pollInterval = time.Second

// WaitReady blocks until every object is serving: Deployments and StatefulSets have all
// their replicas available, LoadBalancer Services have an ingress address, and Pods are
// Ready. Other kinds are taken to be ready as soon as they are applied.
func (c *Client) WaitReady(ctx context.Context, objs []*Object, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for _, obj := range objs {
		var err error

		switch obj.GetKind() {
		case "Deployment":
			err = c.WaitRollout(ctx, obj.GetName(), time.Until(deadline))
		case "StatefulSet":
			err = c.waitStatefulSet(ctx, obj.GetName(), time.Until(deadline))
		case "Service":
			err = c.waitService(ctx, obj.GetName(), time.Until(deadline))
		case "Pod":
			err = c.waitPod(ctx, obj.GetName(), time.Until(deadline))
		default:
			continue
		}

		if err != nil {
			return fmt.Errorf("%s/%s never became ready: %w\n%s",
				obj.GetKind(), obj.GetName(), err, c.Diagnose(ctx, obj))
		}
	}

	return nil
}

func (c *Client) waitStatefulSet(ctx context.Context, name string, timeout time.Duration) error {
	return poll(ctx, timeout, func() (bool, string, error) {
		s, err := c.Clientset.AppsV1().StatefulSets(c.Namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, "not found yet", ignoreNotFound(err)
		}

		want := int32(1)
		if s.Spec.Replicas != nil {
			want = *s.Spec.Replicas
		}

		if s.Status.ReadyReplicas >= want {
			return true, "", nil
		}

		return false, fmt.Sprintf("%d/%d replicas ready", s.Status.ReadyReplicas, want), nil
	})
}

// waitService waits for a LoadBalancer Service to be given an external address. Services
// of any other type need no waiting.
func (c *Client) waitService(ctx context.Context, name string, timeout time.Duration) error {
	return poll(ctx, timeout, func() (bool, string, error) {
		s, err := c.Clientset.CoreV1().Services(c.Namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, "not found yet", ignoreNotFound(err)
		}

		if s.Spec.Type != corev1.ServiceTypeLoadBalancer {
			return true, "", nil
		}

		if len(s.Status.LoadBalancer.Ingress) > 0 {
			return true, "", nil
		}

		return false, "no external address assigned yet", nil
	})
}

func (c *Client) waitPod(ctx context.Context, name string, timeout time.Duration) error {
	return poll(ctx, timeout, func() (bool, string, error) {
		p, err := c.Clientset.CoreV1().Pods(c.Namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, "not found yet", ignoreNotFound(err)
		}

		for _, cond := range p.Status.Conditions {
			if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
				return true, "", nil
			}
		}

		return false, fmt.Sprintf("pod phase %s", p.Status.Phase), nil
	})
}

// LoadBalancerIP returns the external address of a LoadBalancer Service, waiting for one
// to be assigned. It reports the hostname if the provider assigned a name instead of an IP.
func (c *Client) LoadBalancerIP(ctx context.Context, name string, timeout time.Duration) (string, error) {
	var addr string

	err := poll(ctx, timeout, func() (bool, string, error) {
		s, err := c.Clientset.CoreV1().Services(c.Namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, "not found yet", ignoreNotFound(err)
		}

		if s.Spec.Type != corev1.ServiceTypeLoadBalancer {
			return false, "", fmt.Errorf("Service %q has type %s, the lab asks for a LoadBalancer",
				name, s.Spec.Type)
		}

		for _, ing := range s.Status.LoadBalancer.Ingress {
			if ing.IP != "" {
				addr = ing.IP
				return true, "", nil
			}
			if ing.Hostname != "" {
				addr = ing.Hostname
				return true, "", nil
			}
		}

		return false, "no external address assigned yet", nil
	})

	return addr, err
}

// ServicePort returns the port a Service publishes. When the Service has several ports the
// first one wins, which is enough for the labs.
func (c *Client) ServicePort(ctx context.Context, name string) (int32, error) {
	s, err := c.Clientset.CoreV1().Services(c.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return 0, fmt.Errorf("get Service %q: %w", name, err)
	}

	if len(s.Spec.Ports) == 0 {
		return 0, fmt.Errorf("Service %q publishes no ports", name)
	}

	return s.Spec.Ports[0].Port, nil
}

// poll runs check until it reports done, the timeout expires, or it returns an error. The
// message from the last unsuccessful check is folded into the timeout error, so a failure
// says what the cluster was actually stuck on.
func poll(ctx context.Context, timeout time.Duration, check func() (bool, string, error)) error {
	if timeout <= 0 {
		return fmt.Errorf("timed out before the check could start")
	}

	deadline := time.Now().Add(timeout)
	last := "no status yet"

	for {
		done, msg, err := check()
		if err != nil {
			return err
		}
		if done {
			return nil
		}
		if msg != "" {
			last = msg
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timed out after %s: %s", timeout, last)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

func ignoreNotFound(err error) error {
	if apierrors.IsNotFound(err) {
		return nil
	}
	return err
}

// WaitGone blocks until none of the objects exist any more, so that one lab does not
// inherit the previous lab's workloads.
func (c *Client) WaitGone(ctx context.Context, objs []*Object, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for _, obj := range objs {
		ri, err := c.resourceFor(obj)
		if err != nil {
			continue
		}

		name, kind := obj.GetName(), obj.GetKind()
		err = poll(ctx, time.Until(deadline), func() (bool, string, error) {
			_, err := ri.Get(ctx, name, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return true, "", nil
			}
			if err != nil {
				return false, "", err
			}
			return false, "still terminating", nil
		})

		if err != nil {
			return fmt.Errorf("%s/%s was not removed: %w", kind, name, err)
		}
	}

	return nil
}

// PodNames lists the names of the pods matching a label selector.
func (c *Client) PodNames(ctx context.Context, selector string) ([]string, error) {
	pods, err := c.Clientset.CoreV1().Pods(c.Namespace).List(ctx,
		metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, fmt.Errorf("list pods %q: %w", selector, err)
	}

	names := []string{}
	for _, p := range pods.Items {
		names = append(names, p.Name)
	}

	return names, nil
}

// RestartCount returns the total number of container restarts across the pods matching a
// label selector. A liveness probe pointed at the wrong path shows up here.
func (c *Client) RestartCount(ctx context.Context, selector string) (int32, error) {
	pods, err := c.Clientset.CoreV1().Pods(c.Namespace).List(ctx,
		metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return 0, fmt.Errorf("list pods %q: %w", selector, err)
	}

	var total int32
	for _, p := range pods.Items {
		for _, cs := range p.Status.ContainerStatuses {
			total += cs.RestartCount
		}
	}

	return total, nil
}

// WaitPodsReady waits until the given number of pods matching a selector are Ready and no
// others linger, which is what "the rollout settled" means for the lab tests.
func (c *Client) WaitPodsReady(ctx context.Context, selector string, want int, timeout time.Duration) error {
	return poll(ctx, timeout, func() (bool, string, error) {
		pods, err := c.Clientset.CoreV1().Pods(c.Namespace).List(ctx,
			metav1.ListOptions{LabelSelector: selector})
		if err != nil {
			return false, "", err
		}

		ready, other := 0, []string{}
		for _, p := range pods.Items {
			if p.DeletionTimestamp != nil {
				other = append(other, p.Name+"(terminating)")
				continue
			}

			isReady := false
			for _, cond := range p.Status.Conditions {
				if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
					isReady = true
				}
			}

			if isReady {
				ready++
			} else {
				other = append(other, fmt.Sprintf("%s(%s)", p.Name, p.Status.Phase))
			}
		}

		if ready == want && len(other) == 0 {
			return true, "", nil
		}

		return false, fmt.Sprintf("%d/%d pods ready, waiting on %s",
			ready, want, strings.Join(other, " ")), nil
	})
}

// WaitWorkloadsDrained blocks until the objects have really let go of the cluster: the pods
// they owned are gone, and the load balancers of their Services have released their
// addresses.
//
// Deleting an object does not do this synchronously. A Deployment disappears from the API
// while its pods are still terminating, and a Service of type LoadBalancer keeps its
// address until the provider removes the forwarding it set up. Labs share a cluster, so the
// next one has to start from an empty one.
func (c *Client) WaitWorkloadsDrained(ctx context.Context, objs []*Object, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for _, obj := range objs {
		selector := selectorOf(obj)
		if selector == "" {
			continue
		}

		kind, name := obj.GetKind(), obj.GetName()

		switch kind {
		case "Deployment", "StatefulSet":
			err := poll(ctx, time.Until(deadline), func() (bool, string, error) {
				pods, err := c.Clientset.CoreV1().Pods(c.Namespace).List(ctx,
					metav1.ListOptions{LabelSelector: selector})
				if err != nil {
					return false, "", err
				}
				if len(pods.Items) == 0 {
					return true, "", nil
				}

				return false, fmt.Sprintf("%d pod(s) still terminating", len(pods.Items)), nil
			})
			if err != nil {
				return fmt.Errorf("pods of %s/%s did not go away: %w", kind, name, err)
			}

		case "Service":
			// The load balancer of a Service is implemented by its own workload, which
			// the provider removes once the Service is gone.
			err := poll(ctx, time.Until(deadline), func() (bool, string, error) {
				ds, err := c.Clientset.AppsV1().DaemonSets("kube-system").List(ctx,
					metav1.ListOptions{})
				if err != nil {
					// Not every cluster implements load balancers this way.
					return true, "", nil
				}

				for i := range ds.Items {
					if strings.Contains(ds.Items[i].Name, name) {
						return false, "the load balancer is still being torn down", nil
					}
				}

				return true, "", nil
			})
			if err != nil {
				return fmt.Errorf("the load balancer of Service %q was not released: %w", name, err)
			}
		}
	}

	return nil
}
