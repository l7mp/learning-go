//go:build !kvstoremode && !failingkvstoremode

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The local data layer is the default: a test run with no mode tag at all expects the app
// to keep its accounts in memory. The other modes are selected with a build tag, and the
// build constraint above makes sure exactly one mode is ever compiled in.

// TestMain runs before every other test in the package. The local data layer needs nothing
// of its environment, so there is nothing to check here.
func TestMain(m *testing.M) {
	m.Run()
}

// TestDataLayerIsLocal checks that the app keeps its accounts to itself. It can only do so
// when a key-value store is around to be inspected: with no store running there is nothing
// the app could have written to, and nothing to prove.
func TestDataLayerIsLocal(t *testing.T) {
	if !kvstoreReachable() {
		t.Logf("no key-value store at %s, so there is nothing the app could be using "+
			"instead of the local data layer", kvstoreURL())
		return
	}

	assert.Empty(t, transferAndReadStore(t),
		"the key-value store was written to: the app is using the kvstore data layer "+
			"even though this test expects the local one. If you meant to test the "+
			"key-value store, add the kvstoremode build tag")
}
