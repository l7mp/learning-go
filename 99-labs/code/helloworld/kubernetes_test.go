//go:build kubernetes

package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorldKubernetes(t *testing.T) {
	// clean up cluster
	execCmd(t, "kubectl", "delete", "-f", "deploy/kubernetes-deployment.yaml")
	execCmd(t, "kubectl", "delete", "-f", "deploy/kubernetes-service.yaml")

	// redeploy
	execCmd(t, "kubectl", "apply", "-f", "deploy/kubernetes-deployment.yaml")
	execCmd(t, "kubectl", "apply", "-f", "deploy/kubernetes-service.yaml")

	// may take a while
	time.Sleep(10 * time.Second)

	ip, _ := execCmd(t, "kubectl", "get", "service", "helloworld", "-o", `jsonpath="{.status.loadBalancer.ingress[0].ip}"`)
	if ip == "" {
		// may take a while
		time.Sleep(10 * time.Second)
		ip, _ = execCmd(t, "kubectl", "get", "service", "helloworld", "-o", `jsonpath="{.status.loadBalancer.ingress[0].ip}"`)
	}

	ip = strings.Trim(ip, `"`)
	assert.NotNil(t, net.ParseIP(ip), "public IP")

	// and Get
	url := fmt.Sprintf("http://%s//:8080", ip)
	res, err := http.Get(url)
	assert.NoError(t, err, "GET")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status code")

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err, "read response body")

	assert.Regexp(t, regexp.MustCompile("^Hello world from .* running Go version"), string(body), "response")
}
