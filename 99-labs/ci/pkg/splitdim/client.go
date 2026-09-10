// Package splitdim is a client for the SplitDim API under test.
package splitdim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

// Account is the balance of one user.
type Account struct {
	Holder  string `json:"holder"`
	Balance int    `json:"balance"`
}

// Transfer is a money transfer between two users.
type Transfer struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int    `json:"amount"`
}

// Client talks to a SplitDim deployment.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// New returns a client for the SplitDim instance at the given base URL.
func New(baseURL string) *Client {
	return &Client{BaseURL: baseURL, HTTP: httpClient(30 * time.Second)}
}

// WithTimeout returns a copy of the client that gives up after d. Tests use it to put a
// bound on calls that an unfixed app would never answer.
func (c *Client) WithTimeout(d time.Duration) *Client {
	return &Client{BaseURL: c.BaseURL, HTTP: httpClient(d)}
}

// httpClient builds a client that opens a fresh connection per request. A Service spreads
// traffic per connection, not per request, so a client that reused one keep-alive
// connection would talk to a single pod for the whole test and make every check about
// replicas silently vacuous.
func httpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{DisableKeepAlives: true},
	}
}

// Reset zeroes all balances.
func (c *Client) Reset() error {
	status, _, err := c.get("/api/reset")
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("GET /api/reset: HTTP status %d", status)
	}
	return nil
}

// Healthz probes the health check endpoint and returns its HTTP status.
func (c *Client) Healthz() (int, error) {
	status, _, err := c.get("/healthz")
	return status, err
}

// Transfer registers a transfer and returns the HTTP status and how long the call took.
func (c *Client) Transfer(t Transfer) (int, time.Duration, error) {
	body, err := json.Marshal(t)
	if err != nil {
		return 0, 0, fmt.Errorf("marshal transfer: %w", err)
	}

	start := time.Now()
	res, err := c.HTTP.Post(c.BaseURL+"/api/transfer", "application/json", bytes.NewReader(body))
	elapsed := time.Since(start)
	if err != nil {
		return 0, elapsed, fmt.Errorf("POST /api/transfer: %w", err)
	}
	defer res.Body.Close()
	io.Copy(io.Discard, res.Body)

	return res.StatusCode, elapsed, nil
}

// Accounts returns the current balances, sorted by holder so that results from different
// replicas can be compared directly.
func (c *Client) Accounts() ([]Account, error) {
	status, body, err := c.get("/api/accounts")
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GET /api/accounts: HTTP status %d", status)
	}

	accounts := []Account{}
	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, fmt.Errorf("GET /api/accounts: decode %q: %w", string(body), err)
	}

	sort.Slice(accounts, func(i, j int) bool { return accounts[i].Holder < accounts[j].Holder })

	return accounts, nil
}

// Clear returns the transfers that would settle all debts. It doubles as a consistency
// check: the app refuses to clear a ledger whose balances do not add up to zero.
func (c *Client) Clear() ([]Transfer, error) {
	status, body, err := c.get("/api/clear")
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("GET /api/clear: HTTP status %d (the account database is "+
			"most likely inconsistent, balances no longer add up to zero)", status)
	}

	transfers := []Transfer{}
	if err := json.Unmarshal(body, &transfers); err != nil {
		return nil, fmt.Errorf("GET /api/clear: decode %q: %w", string(body), err)
	}

	sort.Slice(transfers, func(i, j int) bool {
		if transfers[i].Sender != transfers[j].Sender {
			return transfers[i].Sender < transfers[j].Sender
		}
		return transfers[i].Receiver < transfers[j].Receiver
	})

	return transfers, nil
}

// TotalBalance sums all balances. A consistent ledger sums to zero.
func (c *Client) TotalBalance() (int, error) {
	accounts, err := c.Accounts()
	if err != nil {
		return 0, err
	}

	total := 0
	for _, a := range accounts {
		total += a.Balance
	}

	return total, nil
}

// readySamples is how many consecutive successful probes make a deployment "ready".
//
// One success is not enough. After a rollout the load balancer keeps forwarding to the
// endpoint it knew about for a short while, so a probe can succeed against the pod that is
// on its way out and the very next request is refused. Insisting on a few consecutive
// successes waits for the data plane to catch up with the control plane.
const readySamples = 3

// WaitReady blocks until the app answers reliably through its Service, so that tests do
// not race either a pod that is still starting or a load balancer that has not yet been
// reprogrammed.
func (c *Client) WaitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	probe := c.WithTimeout(5 * time.Second)

	var last error
	streak := 0

	for time.Now().Before(deadline) {
		status, _, err := probe.get("/api/accounts")
		switch {
		case err != nil:
			last, streak = err, 0
		case status != http.StatusOK:
			last, streak = fmt.Errorf("HTTP status %d", status), 0
		default:
			streak++
			if streak >= readySamples {
				return nil
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("SplitDim at %s did not become ready in %s: %w", c.BaseURL, timeout, last)
}

func (c *Client) get(path string) (int, []byte, error) {
	res, err := c.HTTP.Get(c.BaseURL + path)
	if err != nil {
		return 0, nil, fmt.Errorf("GET %s: %w", path, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, fmt.Errorf("GET %s: read body: %w", path, err)
	}

	return res.StatusCode, body, nil
}
