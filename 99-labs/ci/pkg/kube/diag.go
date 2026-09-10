package kube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// maxLogLines caps how much of a container log we quote into a failure message.
const maxLogLines = 40

// Diagnose collects the events and container logs that explain why an object is not
// coming up. A student reading a CI failure should not have to reproduce it locally to
// find out what went wrong.
func (c *Client) Diagnose(ctx context.Context, obj *Object) string {
	var b strings.Builder

	fmt.Fprintf(&b, "\n--- diagnostics for %s/%s ---\n", obj.GetKind(), obj.GetName())

	if ev := c.events(ctx, obj); ev != "" {
		fmt.Fprintf(&b, "events:\n%s", ev)
	}

	selector := selectorOf(obj)
	if selector == "" {
		return b.String()
	}

	pods, err := c.Clientset.CoreV1().Pods(c.Namespace).List(ctx,
		metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		fmt.Fprintf(&b, "could not list pods (%s): %s\n", selector, err)
		return b.String()
	}

	for i := range pods.Items {
		p := &pods.Items[i]
		fmt.Fprintf(&b, "\npod %s: phase %s\n", p.Name, p.Status.Phase)

		for _, cs := range p.Status.ContainerStatuses {
			fmt.Fprintf(&b, "  container %s: ready=%t restarts=%d%s\n",
				cs.Name, cs.Ready, cs.RestartCount, waitingReason(cs))
		}

		for _, cond := range p.Status.Conditions {
			if cond.Status != corev1.ConditionTrue && cond.Message != "" {
				fmt.Fprintf(&b, "  %s: %s\n", cond.Type, cond.Message)
			}
		}

		if logs := c.podLogs(ctx, p.Name); logs != "" {
			fmt.Fprintf(&b, "  logs (last %d lines):\n%s", maxLogLines, indent(logs))
		}
	}

	return b.String()
}

func waitingReason(cs corev1.ContainerStatus) string {
	if cs.State.Waiting != nil && cs.State.Waiting.Reason != "" {
		msg := cs.State.Waiting.Reason
		if cs.State.Waiting.Message != "" {
			msg += ": " + cs.State.Waiting.Message
		}
		return " waiting=" + msg
	}
	if cs.State.Terminated != nil {
		return fmt.Sprintf(" terminated=%s(exit %d)",
			cs.State.Terminated.Reason, cs.State.Terminated.ExitCode)
	}
	return ""
}

// events returns the recent events for an object, newest last.
func (c *Client) events(ctx context.Context, obj *Object) string {
	list, err := c.Clientset.CoreV1().Events(c.Namespace).List(ctx, metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + obj.GetName(),
	})
	if err != nil || len(list.Items) == 0 {
		return ""
	}

	sort.Slice(list.Items, func(i, j int) bool {
		return list.Items[i].LastTimestamp.Before(&list.Items[j].LastTimestamp)
	})

	var b strings.Builder
	for _, e := range list.Items {
		fmt.Fprintf(&b, "  %s %s: %s\n", e.Type, e.Reason, e.Message)
	}

	return b.String()
}

// podLogs returns the tail of a pod's log, including the previous container instance when
// the pod has been restarting.
func (c *Client) podLogs(ctx context.Context, name string) string {
	lines := int64(maxLogLines)

	req := c.Clientset.CoreV1().Pods(c.Namespace).GetLogs(name, &corev1.PodLogOptions{
		TailLines: &lines,
	})

	stream, err := req.Stream(ctx)
	if err != nil {
		return ""
	}
	defer stream.Close()

	body, err := io.ReadAll(stream)
	if err != nil {
		return ""
	}

	return string(body)
}

// selectorOf builds a label selector that matches the pods belonging to an object.
func selectorOf(obj *Object) string {
	labels, found, err := nestedStringMap(obj, "spec", "selector", "matchLabels")
	if err != nil || !found {
		// Services carry a flat selector rather than a matchLabels block.
		labels, found, err = nestedStringMap(obj, "spec", "selector")
		if err != nil || !found {
			return ""
		}
	}

	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+labels[k])
	}

	return strings.Join(parts, ",")
}

func nestedStringMap(obj *Object, fields ...string) (map[string]string, bool, error) {
	cur := obj.Object
	for _, f := range fields {
		next, ok := cur[f]
		if !ok {
			return nil, false, nil
		}
		m, ok := next.(map[string]any)
		if !ok {
			return nil, false, nil
		}
		cur = m
	}

	ret := map[string]string{}
	for k, v := range cur {
		s, ok := v.(string)
		if !ok {
			return nil, false, nil
		}
		ret[k] = s
	}

	return ret, len(ret) > 0, nil
}

func indent(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lines {
		lines[i] = "    " + l
	}
	return strings.Join(lines, "\n") + "\n"
}

// marshal is json.Marshal, kept here so the patch helpers read cleanly.
func marshal(v any) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal patch: %w", err)
	}
	return body, nil
}
