//go:build main

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testServerAvailable(t *testing.T, ctx context.Context) {
	if _, err := http.Get("http://:8080"); err != nil {
		// Server is not running
		go execCmdContext(ctx, t, "go", "run", "main.go")

		// Wait for the server to start
		time.Sleep(500 * time.Millisecond)
	}
}

func TestHelloWorldLocal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testServerAvailable(t, ctx)

	h, err := os.Hostname()
	assert.NoError(t, err, "reading hostname")

	res, err := http.Get("http://:8080")
	assert.NoError(t, err, "GET")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status code")

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err, "read response body")
	// assert.Equal(t, fmt.Sprintf("Hello world from %s!", h), string(body), "response")
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("^Hello world from %s", h)), string(body), "response")
}
