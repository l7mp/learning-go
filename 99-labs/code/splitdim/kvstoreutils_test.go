package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// The helpers below talk to the key-value store directly, over plain HTTP, so that tests
// can check where the account data actually ends up. They deliberately do not import the
// kvstore packages: the SplitDim module does not depend on the key-value store until lab
// 04, and these helpers must keep compiling before that.

// kvEntry is one key-value pair in the store, matching the JSON of the kvstore API.
type kvEntry struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Version int    `json:"version"`
}

// faultConfig configures fault injection on one key-value store API path. The key-value
// store serves this API only when started with KVSTORE_FAULT_INJECTION set.
type faultConfig struct {
	Path         string `json:"path"`
	FailNext     int    `json:"failNext"`
	AbortPercent int    `json:"abortPercent"`
	Status       int    `json:"status"`
	Delay        string `json:"delay"`
}

// faultStatus reports the fault configuration and the request counters of the key-value
// store. The counters let a test check that a client actually retried a failed call.
type faultStatus struct {
	Faults   map[string]faultConfig `json:"faults"`
	Requests map[string]int         `json:"requests"`
	Failures map[string]int         `json:"failures"`
}

// kvstoreURL returns the base URL of the key-value store, localhost:8081 unless the
// KVSTORE_URL environment variable says otherwise.
func kvstoreURL() string {
	if u := os.Getenv("KVSTORE_URL"); u != "" {
		return u
	}
	return "http://localhost:8081"
}

// kvstoreReachable reports whether the key-value store answers at all.
func kvstoreReachable() bool {
	c := http.Client{Timeout: 2 * time.Second}
	res, err := c.Get(kvstoreURL() + "/api/list")
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return res.StatusCode == http.StatusOK
}

// kvstoreList returns everything currently stored in the key-value store.
func kvstoreList(t *testing.T) []kvEntry {
	t.Helper()

	c := http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(kvstoreURL() + "/api/list")
	if err != nil {
		t.Fatalf("key-value store list: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("key-value store list: HTTP status %d", res.StatusCode)
	}

	entries := []kvEntry{}
	if err := json.NewDecoder(res.Body).Decode(&entries); err != nil {
		t.Fatalf("key-value store list: decode: %s", err)
	}

	return entries
}

// kvstoreBalances returns the account balances held in the key-value store.
func kvstoreBalances(t *testing.T) map[string]string {
	t.Helper()

	ret := map[string]string{}
	for _, e := range kvstoreList(t) {
		ret[e.Key] = e.Value
	}

	return ret
}

// kvstoreReset empties the key-value store.
func kvstoreReset(t *testing.T) {
	t.Helper()

	c := http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(kvstoreURL() + "/api/reset")
	if err != nil {
		t.Fatalf("key-value store reset: %s", err)
	}
	defer res.Body.Close()
}

// requireFaultInjection fails the calling test unless the key-value store was started with
// fault injection enabled.
func requireFaultInjection(t *testing.T) {
	t.Helper()

	c := http.Client{Timeout: 5 * time.Second}
	res, err := c.Get(kvstoreURL() + "/api/fault")
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("the key-value store at %s does not serve the fault injection API: "+
			"restart it with KVSTORE_FAULT_INJECTION=1, see the lab README", kvstoreURL())
	}
	res.Body.Close()
}

// requireKVStore stops the whole run when no key-value store answers. It is called from
// TestMain, so a suite that expects the store can never quietly test something else.
func requireKVStore() {
	if kvstoreReachable() {
		return
	}

	fmt.Fprintf(os.Stderr, "\nthis run expects the key-value store to be running at %s, "+
		"start it with:\n    cd 99-labs/code/kvstore && go run kvstore.go\n\n", kvstoreURL())
	os.Exit(1)
}

// refuseKVStore stops the whole run when a key-value store answers, for the checks that are
// about the app coping with a downstream that is not there.
func refuseKVStore() {
	if !kvstoreReachable() {
		return
	}

	fmt.Fprintf(os.Stderr, "\nthis run expects the key-value store at %s to be STOPPED: "+
		"it checks that the app gives up on an unreachable downstream instead of\n"+
		"retrying forever, which proves nothing while the store is healthy\n\n",
		kvstoreURL())
	os.Exit(1)
}

// transferAndReadStore empties both stores, registers one transfer through the SplitDim
// API, and returns what ended up in the key-value store. The mode checks use it to tell
// the data layers apart: where the balances land is the only thing that distinguishes them.
func transferAndReadStore(t *testing.T) map[string]string {
	t.Helper()

	kvstoreReset(t)

	res, err := testHTTP(t, "api/reset", "GET", "")
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("GET api/reset failed, cannot set up the test: %v", err)
	}

	res, err = testHTTP(t, "api/transfer", "POST", `{"sender":"alice","receiver":"bob","amount":7}`)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("POST api/transfer failed, cannot set up the test: %v", err)
	}

	return kvstoreBalances(t)
}

// kvstoreFault installs a fault on a key-value store API path.
func kvstoreFault(t *testing.T, c faultConfig) {
	t.Helper()

	body, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("fault config marshal: %s", err)
	}

	cl := http.Client{Timeout: 10 * time.Second}
	res, err := cl.Post(kvstoreURL()+"/api/fault", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("install fault: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(res.Body)
		t.Fatalf("install fault: HTTP status %d: %s", res.StatusCode, string(msg))
	}
}

// kvstoreFaultClear removes all faults and zeroes the request counters.
func kvstoreFaultClear(t *testing.T) {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, kvstoreURL()+"/api/fault", nil)
	if err != nil {
		t.Fatalf("clear faults: %s", err)
	}

	cl := http.Client{Timeout: 10 * time.Second}
	res, err := cl.Do(req)
	if err != nil {
		t.Fatalf("clear faults: %s", err)
	}
	defer res.Body.Close()
}

// kvstoreFaultStatus returns the fault configuration and request counters.
func kvstoreFaultStatus(t *testing.T) faultStatus {
	t.Helper()

	cl := http.Client{Timeout: 10 * time.Second}
	res, err := cl.Get(kvstoreURL() + "/api/fault")
	if err != nil {
		t.Fatalf("fault status: %s", err)
	}
	defer res.Body.Close()

	s := faultStatus{}
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		t.Fatalf("fault status: decode: %s", err)
	}

	return s
}

// transferWithin posts a transfer and gives up after the deadline. The plain testHTTP
// helper uses a client without a timeout, which would hang forever against an app that
// retries a failing downstream call in an infinite loop.
func transferWithin(t *testing.T, body string, timeout time.Duration) (*http.Response, time.Duration, error) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, testURL("api/transfer"),
		bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("create transfer request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	c := http.Client{Timeout: timeout}

	start := time.Now()
	res, err := c.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		return nil, elapsed, fmt.Errorf("transfer did not complete within %s: %w", timeout, err)
	}

	return res, elapsed, nil
}
