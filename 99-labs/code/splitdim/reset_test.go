//go:build reset

package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test through the HTTP layer.
func TestReset(t *testing.T) {
	// reset!
	res, err := testHTTP(t, "api/reset", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")
}
