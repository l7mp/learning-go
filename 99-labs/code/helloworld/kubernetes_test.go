//go:build kubernetes

package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// clean up cluster
	execCmd(t, "kubectl", "delete", "-f", "deploy/kubernetes-deployment.yaml")
	execCmd(t, "kubectl", "delete", "-f", "deploy/kubernetes-service.yaml")

	// build the container image
	execCmd(t, "minikube", "image", "build", "-t", "helloworld", "-f", "deploy/Dockerfile", ".")
	go execCmd(t, "minikube", "tunnel")

	// redeploy
	execCmd(t, "kubectl", "apply", "-f", "deploy/kubernetes-deployment.yaml")
	execCmd(t, "kubectl", "apply", "-f", "deploy/kubernetes-service.yaml")

	// may take a while
	time.Sleep(10 * time.Second)

	ip, _ := execCmd(t, "kubectl", "get", "service", "helloworld", "-o", `jsonpath="{.status.loadBalancer.ingress[0].ip}"`)
	if ip == "" {
		// make sure minikube tunnel is running if no public IP exists
		execCmdContext(ctx, t, "minikube", "tunnel")
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
