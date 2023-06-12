//go:build httphandler

package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type endpoint struct {
	api, method string
	status      int
}

// TestAPIEndpoints checks whether the HTTP server is available, the empty path for the
// static HTML and all API endpoints are accessible.
func TestAPIEndpoints(t *testing.T) {
	testEndpoints := []endpoint{
		{"", "GET", http.StatusOK},
		{"api/transfer", "POST", http.StatusOK},
		{"api/accounts", "GET", http.StatusOK},
		{"api/clear", "GET", http.StatusOK},
		{"api/reset", "GET", http.StatusOK},
	}

	for _, ep := range testEndpoints {
		res, err := testHTTP(t, ep.api, ep.method, "")
		assert.NoError(t, err, fmt.Sprintf("%s: %s should return status", ep.method, ep.api))
		assert.Equal(t, ep.status, res.StatusCode,
			fmt.Sprintf("%s: %s status err: expected: %d, got: %d",
				ep.method, ep.api, ep.status, res.StatusCode))
	}
}

// TestAPIMethods checks whether the HTTP server is available, the empty path for the
// static HTML and all API endpoints are accessible, and only the permitted HTTP methods are
// implemented.
func TestAPIMethods(t *testing.T) {
	testEndpoints := []endpoint{
		{"", "GET", http.StatusOK},
		{"api/transfer", "POST", http.StatusOK},
		{"api/transfer", "GET", http.StatusMethodNotAllowed},
		{"api/accounts", "GET", http.StatusOK},
		{"api/accounts", "POST", http.StatusMethodNotAllowed},
		{"api/clear", "GET", http.StatusOK},
		{"api/clear", "POST", http.StatusMethodNotAllowed},
		{"api/reset", "GET", http.StatusOK},
		{"api/reset", "POST", http.StatusMethodNotAllowed},
	}

	for _, ep := range testEndpoints {
		res, err := testHTTP(t, ep.api, ep.method, "")
		assert.NoError(t, err, fmt.Sprintf("%s: %s should return status", ep.method, ep.api))
		assert.Equal(t, ep.status, res.StatusCode,
			fmt.Sprintf("%s: %s status err: expected: %d, got: %d",
				ep.method, ep.api, ep.status, res.StatusCode))
	}
}

// // TestAPIPaths checks whether the HTTP server is available, the empty path for the
// // static HTML and all API endpoints are accessible, and only the permitted HTTP methods are
// // implemented and only on the requested paths.
// func TestAPIPaths(t *testing.T) {
// 	testEndpoints := []struct {
// 		api, method string
// 		status      int
// 	}{
// 		{"", "GET", http.StatusOK},
// 		{"/", "GET", http.StatusOK},
// 		{"api/transfer", "POST", http.StatusOK},
// 		{"api/transfer", "GET", http.StatusMethodNotAllowed},
// 		{"api/transfer/dummy", "POST", http.StatusNotFound},
// 		{"api/accounts", "GET", http.StatusOK},
// 		{"api/accounts", "POST", http.StatusMethodNotAllowed},
// 		{"api/accounts/dummy", "GET", http.StatusNotFound},
// 		{"api/clear", "GET", http.StatusOK},
// 		{"api/clear", "POST", http.StatusMethodNotAllowed},
// 		{"api/clear/dummy", "GET", http.StatusNotFound},
// 		{"api/reset", "GET", http.StatusOK},
// 		{"api/reset", "POST", http.StatusMethodNotAllowed},
// 		{"api/reset/dummy", "GET", http.StatusNotFound},
// 	}

// 	for _, ep := range testEndpoints {
//		res, err := testHTTP(t, ep.api, ep.method, "")
// 		assert.NoError(t, err, fmt.Sprintf("%s: %s should fail", ep.method, ep.api))
// 		assert.Equal(t, ep.status, res.StatusCode,
// 			fmt.Sprintf("%s: %s status err: expected: %d, got: %d",
// 				ep.method, ep.api, ep.status, res.StatusCode))
// 	}
// }
