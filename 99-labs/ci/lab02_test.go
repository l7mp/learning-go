package ci

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"labci/pkg/labenv"
)

// Lab 02 asks the student to build a "hello world" web app, package it into a container,
// and expose it from Kubernetes behind a LoadBalancer Service.

// TestLab02 deploys the helloworld app and checks it end to end.
func TestLab02(t *testing.T) {
	e := labenv.SetupHelloWorld(t)
	defer e.Teardown()

	base := fmt.Sprintf("http://%s:%d", e.HelloIP, e.HelloPort)

	t.Run("GoVersionSuite", func(t *testing.T) {
		// The student's own test, run against the Kubernetes deployment.
		e.RunSuite(t, "helloworld", "helloworldgoversion")
	})

	t.Run("Greeting", func(t *testing.T) {
		body, status := get(t, base+"/")
		assert.Equal(t, http.StatusOK, status, "GET / through the LoadBalancer")
		assert.Regexp(t, regexp.MustCompile(`^Hello world from .* running Go version .*`), body,
			"the greeting should name both the hostname and the Go version")
	})

	// A Deployment behind a Service spreads requests over its pods. A student who
	// deployed a bare Pod, or who wired the Service selector to a single pod, passes
	// every other check but fails this one.
	t.Run("ServiceLoadBalancesOverReplicas", func(t *testing.T) {
		if err := e.Kube.Scale(e.Ctx, "helloworld", 3); err != nil {
			t.Fatalf("%s", err)
		}
		if err := e.Kube.WaitRollout(e.Ctx, "helloworld", 3*time.Minute); err != nil {
			t.Fatalf("scaling the helloworld Deployment to 3 replicas: %s", err)
		}
		waitServing(t, base+"/")

		seen := map[string]int{}
		hostname := regexp.MustCompile(`^Hello world from ([^ ]+) `)

		for i := 0; i < 40; i++ {
			body, status := get(t, base+"/")
			if status != http.StatusOK {
				t.Fatalf("GET / returned HTTP %d", status)
			}

			m := hostname.FindStringSubmatch(body)
			if m == nil {
				t.Fatalf("cannot find a hostname in the greeting %q", body)
			}
			seen[m[1]]++
		}

		assert.GreaterOrEqual(t, len(seen), 2,
			"40 requests all landed on the same pod (%v): the Service does not appear to "+
				"be spreading traffic over the replicas of a Deployment", seen)
	})
}

// waitServing blocks until the URL answers reliably. A completed rollout does not yet
// mean the Service forwards to it: the load balancer is reprogrammed a moment later, so a
// single successful request is not proof and the next one can still be refused.
func waitServing(t *testing.T, url string) {
	t.Helper()

	const samples = 3

	c := http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{DisableKeepAlives: true},
	}

	deadline := time.Now().Add(2 * time.Minute)
	streak := 0

	for time.Now().Before(deadline) {
		res, err := c.Get(url)
		if err == nil {
			res.Body.Close()
		}

		if err == nil && res.StatusCode == http.StatusOK {
			streak++
			if streak >= samples {
				return
			}
		} else {
			streak = 0
		}

		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf("%s did not start serving reliably", url)
}

// get issues a GET and returns the body and the status code.
func get(t *testing.T, url string) (string, int) {
	t.Helper()

	// A fresh connection per request: a Service load balances per connection, so a
	// client that kept one alive would only ever reach a single pod.
	c := http.Client{
		Timeout:   30 * time.Second,
		Transport: &http.Transport{DisableKeepAlives: true},
	}

	res, err := c.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %s", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("GET %s: read body: %s", url, err)
	}

	return string(body), res.StatusCode
}
