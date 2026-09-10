//go:build kvstoremode

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMain runs before every other test in the package, so a run tagged kvstoremode can
// never test the wrong thing: without a key-value store the whole run stops here.
func TestMain(m *testing.M) {
	requireKVStore()
	m.Run()
}

// TestDataLayerIsKVStore checks that the accounts really live in the key-value store. The
// rest of the suite passes just as happily against the in-memory data layer, because the
// SplitDim API behaves identically either way, so this is the test that tells them apart.
func TestDataLayerIsKVStore(t *testing.T) {
	balances := transferAndReadStore(t)

	assert.Equal(t, "7", balances["alice"],
		"the key-value store does not hold alice's balance: the app is not using the "+
			"kvstore data layer")
	assert.Equal(t, "-7", balances["bob"],
		"the key-value store does not hold bob's balance")
}
