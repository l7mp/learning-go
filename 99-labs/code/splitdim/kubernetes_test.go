//go:build kubernetes

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSplitDimKubernetes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// clean up cluster
	execCmd(t, "kubectl", "delete", "-f", "deploy/kubernetes-local-db.yaml")

	// build the container image
	_, _, err := execCmd(t, "minikube", "image", "build", "-t", "splitdim", "-f", "deploy/Dockerfile", ".")
	assert.NoError(t, err, "kubectl delete")

	// redeploy
	_, _, err = execCmd(t, "kubectl", "apply", "-f", "deploy/kubernetes-local-db.yaml")
	assert.NoError(t, err, "kubectl apply")

	// may take a while
	time.Sleep(20 * time.Second)

	// redeploy
	_, _, err = execCmd(t, "kubectl", "get", "service", "splitdim")
	assert.NoError(t, err, "kubectl get")

	ip, _, err := execCmd(t, "kubectl", "get", "service", "splitdim", "-o", `jsonpath="{.status.loadBalancer.ingress[0].ip}"`)
	if ip == "" {
		// make sure minikube tunnel is running if no public IP exists
		execCmdContext(ctx, t, "minikube", "tunnel")
		// may take a while
		time.Sleep(10 * time.Second)
		ip, _, err = execCmd(t, "kubectl", "get", "service", "splitdim", "-o", `jsonpath="{.status.loadBalancer.ingress[0].ip}"`)
	}

	ip = strings.Trim(ip, `"`)
	assert.NotNil(t, net.ParseIP(ip), "public IP")

	// reset
	url := fmt.Sprintf("http://%s:80/api/reset", ip)
	res, err := http.Get(url)
	assert.NoError(t, err, "GET")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status code")

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err, "read response body")
	assert.Equal(t, "", string(body), "response")

	// transfer
	url = fmt.Sprintf("http://%s:80/api/transfer", ip)
	b := bytes.NewBuffer([]byte(`{"sender":"a","receiver":"b","amount":1}`))
	res, err = http.Post(url, "application/json", b)
	assert.NoError(t, err, "POST")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status code")

	// accounts
	url = fmt.Sprintf("http://%s:80/api/accounts", ip)
	res, err = http.Get(url)
	assert.NoError(t, err, "GET")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status code")

	body, err = io.ReadAll(res.Body)
	assert.NoError(t, err, "read response body")
	assert.Equal(t, `[{"holder":"a","balance":1},{"holder":"b","balance":-1}]`, string(body), "response")

	// clear
	url = fmt.Sprintf("http://%s:80/api/clear", ip)
	res, err = http.Get(url)
	assert.NoError(t, err, "GET")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status code")

	body, err = io.ReadAll(res.Body)
	assert.NoError(t, err, "read response body")
	assert.Equal(t, `[{"sender":"b","receiver":"a","amount":1}]`, string(body), "response")

	// clean up cluster
	execCmd(t, "kubectl", "delete", "-f", "deploy/kubernetes-local-db.yaml")

}
