//go:build retry

package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// budgetDeadline bounds how long a transfer may take. The retry policy the lab suggests
// finishes in about a second.
const budgetDeadline = 5 * time.Second

// TestRetryBudget checks that the retry policy is bounded but larger than a single
// attempt: a handful of transient failures must be absorbed, a permanent one must not.
// Start the key-value store with KVSTORE_FAULT_INJECTION=1 before running it.
func TestRetryBudget(t *testing.T) {
	requireFaultInjection(t)

	kvstoreFaultClear(t)
	defer kvstoreFaultClear(t)

	res, err := testHTTP(t, "api/reset", "GET", "")
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("GET api/reset failed, cannot set up the test: %v", err)
	}

	// A few transient failures must be absorbed by the retry policy.
	kvstoreFault(t, faultConfig{Path: "/api/put", FailNext: 3})

	res, _, err = transferWithin(t, `{"sender":"a","receiver":"b","amount":1}`, budgetDeadline)
	if err != nil {
		t.Fatalf("the transfer never returned: %s", err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode,
		"the transfer failed on 3 transient key-value store errors, the retry policy "+
			"should have absorbed them")

	status := kvstoreFaultStatus(t)
	assert.Greater(t, status.Requests["/api/put"], 1,
		"SplitDim issued a single put and gave up: it did not retry at all")

	// A permanent failure must not be retried forever.
	kvstoreFaultClear(t)
	kvstoreFault(t, faultConfig{Path: "/api/put", FailNext: 1000})

	res, elapsed, err := transferWithin(t, `{"sender":"c","receiver":"d","amount":1}`,
		budgetDeadline)
	if err != nil {
		t.Fatalf("the transfer never returned, the retry policy does not appear to be "+
			"bounded: %s", err)
	}

	assert.NotEqual(t, http.StatusOK, res.StatusCode,
		"the transfer reported success even though every put failed")
	assert.Less(t, elapsed, budgetDeadline, "the transfer took too long to give up")
}
