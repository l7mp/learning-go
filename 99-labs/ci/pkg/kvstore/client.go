// Package kvstore drives the key-value store's fault injection, which is a CI-only API
// that a student build does not serve.
//
// It deliberately stops there. Reading and clearing the store are plain HTTP calls in the
// tests instead, because lab 04 asks students to write exactly that client library and a
// finished one in this repository would be the answer to their exercise.
package kvstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FaultConfig configures fault injection on one API path. The store serves this API only
// when it was started with KVSTORE_FAULT_INJECTION set, which the CI does and a student
// build does not.
type FaultConfig struct {
	Path         string `json:"path"`
	FailNext     int    `json:"failNext"`
	AbortPercent int    `json:"abortPercent"`
	Status       int    `json:"status"`
	Delay        string `json:"delay"`
}

// FaultStatus reports the configured faults and the request counters.
type FaultStatus struct {
	Faults   map[string]FaultConfig `json:"faults"`
	Requests map[string]int         `json:"requests"`
	Failures map[string]int         `json:"failures"`
}

// Client talks to a key-value store deployment.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// New returns a client for the key-value store at the given base URL.
func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout:   30 * time.Second,
			Transport: &http.Transport{DisableKeepAlives: true},
		},
	}
}

// Reachable reports whether the store answers at all.
func (c *Client) Reachable() bool {
	probe := http.Client{Timeout: 3 * time.Second}

	res, err := probe.Get(c.BaseURL + "/api/list")
	if err != nil {
		return false
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK
}

// readySamples is how many consecutive successful probes make the store "ready". One is
// not enough: while a replacement pod is coming up the load balancer can still be
// forwarding to the one that is on its way out, so a probe succeeds and the next request
// is refused.
const readySamples = 3

// WaitReady blocks until the store answers reliably.
func (c *Client) WaitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	streak := 0

	for time.Now().Before(deadline) {
		if c.Reachable() {
			streak++
			if streak >= readySamples {
				return nil
			}
		} else {
			streak = 0
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("key-value store at %s did not become ready in %s", c.BaseURL, timeout)
}

// Fault installs a fault on an API path.
func (c *Client) Fault(cfg FaultConfig) error {
	body, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal fault config: %w", err)
	}

	res, err := c.HTTP.Post(c.BaseURL+"/api/fault", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("POST /api/fault: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(res.Body)
		return fmt.Errorf("POST /api/fault: HTTP status %d: %s", res.StatusCode, string(msg))
	}

	return nil
}

// ClearFaults removes every fault and zeroes the request counters.
func (c *Client) ClearFaults() error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL+"/api/fault", nil)
	if err != nil {
		return fmt.Errorf("build fault clear request: %w", err)
	}

	res, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE /api/fault: %w", err)
	}
	defer res.Body.Close()
	io.Copy(io.Discard, res.Body)

	return nil
}

// FaultStatus returns the configured faults and the request counters.
func (c *Client) FaultStatus() (FaultStatus, error) {
	s := FaultStatus{}

	res, err := c.HTTP.Get(c.BaseURL + "/api/fault")
	if err != nil {
		return s, fmt.Errorf("GET /api/fault: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return s, fmt.Errorf("GET /api/fault: HTTP status %d (the key-value store was not "+
			"started with fault injection enabled)", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		return s, fmt.Errorf("GET /api/fault: decode: %w", err)
	}

	return s, nil
}

// FaultInjectionAvailable reports whether the store serves the fault control API.
func (c *Client) FaultInjectionAvailable() bool {
	_, err := c.FaultStatus()
	return err == nil
}
