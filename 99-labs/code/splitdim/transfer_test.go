//go:build transfer

package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test through the HTTP layer. Actual results will be tested with the accountlist api is also available.
func TestTransfer(t *testing.T) {
	res, err := testHTTP(t, "api/reset", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/transfer", "POST", "")
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	// test with random bllst
	res, err = testHTTP(t, "api/transfer", "POST", "dummy")
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "status")

	// test with wrong JSON
	res, err = testHTTP(t, "api/transfer", "POST", `{"a":12}`)
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "status")

	res, err = testHTTP(t, "api/reset", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")
}
