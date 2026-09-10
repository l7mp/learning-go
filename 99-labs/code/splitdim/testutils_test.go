package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testURL builds the URL of an API endpoint on the SplitDim web service under test. The
// address defaults to localhost:8080 and can be overridden with the EXTERNAL_IP and
// EXTERNAL_PORT environment variables, which is how the tests are pointed at a Kubernetes
// LoadBalancer.
func testURL(api string) string {
	addr := "localhost"
	if os.Getenv("EXTERNAL_IP") != "" {
		addr = os.Getenv("EXTERNAL_IP")
	}
	port := "8080"
	if os.Getenv("EXTERNAL_PORT") != "" {
		port = os.Getenv("EXTERNAL_PORT")
	}

	return fmt.Sprintf("http://%s:%s/%s", addr, port, api)
}

// testHTTP issues a request against the SplitDim web service under test.
func testHTTP(t *testing.T, api, method, body string) (*http.Response, error) {
	uri := testURL(api)

	var req *http.Request
	var err error
	if method == "POST" {
		var b *bytes.Buffer
		if body == "" {
			b = bytes.NewBuffer([]byte(`{"sender":"c","receiver":"a", "amount": 4}`))
		} else {
			b = bytes.NewBuffer([]byte(body))
		}
		req, err = http.NewRequest(method, uri, b)
	} else {
		req, err = http.NewRequest(method, uri, nil)
	}

	assert.NoError(t, err, "create req")

	return http.DefaultClient.Do(req)
}
