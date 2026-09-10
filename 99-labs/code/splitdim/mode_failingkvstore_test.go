//go:build failingkvstoremode

package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// hangGuard bounds how long we wait for a transfer, so that an app which retries forever
// reports something useful instead of hanging until the whole test run times out.
//
// It is deliberately generous and is not a performance criterion. How long a bounded retry
// policy takes depends on how the downstream fails: a store whose pods are gone takes its
// Service DNS with them, and each attempt then costs a resolver timeout rather than the
// few hundred milliseconds the backoff asks for. What separates a correct app from a
// broken one here is that it gives up at all, not how quickly.
const hangGuard = 2 * time.Minute

// TestMain runs before every other test in the package. This mode is about the app facing
// a key-value store that is not there, so a store that answers means the check would prove
// nothing and the run stops.
func TestMain(m *testing.M) {
	refuseKVStore()
	m.Run()
}

// TestTransferGivesUp checks that a transfer fails in bounded time when the key-value
// store is unreachable, instead of retrying forever.
//
// This is the one check that separates the code from the previous lab from the code from
// this one: while the store is healthy an unbounded retry loop and a proper retry policy
// behave identically, and only a downstream that never recovers tells them apart.
func TestTransferGivesUp(t *testing.T) {
	res, elapsed, err := transferWithin(t, `{"sender":"a","receiver":"b","amount":1}`,
		hangGuard)
	if err != nil {
		t.Fatalf("the transfer never returned, the app is most likely retrying the "+
			"failing key-value store call in an infinite loop: %s", err)
	}

	t.Logf("the transfer gave up after %s", elapsed)

	assert.NotEqual(t, http.StatusOK, res.StatusCode,
		"the transfer reported success even though the key-value store is down")
}
