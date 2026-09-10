//go:build healthz

package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHealthz checks that the app serves a health check endpoint that Kubernetes can
// probe.
func TestHealthz(t *testing.T) {
	res, err := testHTTP(t, "healthz", "GET", "")
	assert.NoError(t, err, "GET: healthz")
	assert.Equal(t, http.StatusOK, res.StatusCode,
		"GET /healthz should return 200 so that Kubernetes can use it as a liveness probe")
}
