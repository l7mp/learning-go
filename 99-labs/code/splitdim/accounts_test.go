//go:build accounts

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"splitdim/pkg/api"
)

var _ = fmt.Sprintf("dummy")

// Test through the HTTP layer. Some semi-random transfers and the correponnding account lists.
func TestAccounts(t *testing.T) {
	// reset!
	res, err := testHTTP(t, "api/reset", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/transfer", "POST", `{"sender":"a", "receiver":"b", "amount": 4}`)
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/accounts", "GET", "")
	assert.NoError(t, err, "GET: api/accounts")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	a := []api.Account{}
	err = json.NewDecoder(res.Body).Decode(&a)
	assert.NoError(t, err, "account list unmarshal")
	assert.Equal(t, a, []api.Account{{"a", 4}, {"b", -4}}, "a")

	res, err = testHTTP(t, "api/transfer", "POST", `{"sender":"a", "receiver":"b", "amount": 4}`)
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/accounts", "GET", "")
	assert.NoError(t, err, "GET: api/accounts")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	a = []api.Account{}
	err = json.NewDecoder(res.Body).Decode(&a)
	assert.NoError(t, err, "account list unmarshal")
	assert.Equal(t, a, []api.Account{{"a", 8}, {"b", -8}}, "a")

	res, err = testHTTP(t, "api/transfer", "POST", `{"sender":"b", "receiver":"c", "amount": 2}`)
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/accounts", "GET", "")
	assert.NoError(t, err, "GET: api/accounts")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	a = []api.Account{}
	err = json.NewDecoder(res.Body).Decode(&a)
	assert.NoError(t, err, "account list unmarshal")
	assert.Equal(t, a, []api.Account{{"a", 8}, {"b", -6}, {"c", -2}}, "a")

	// clear
	res, err = testHTTP(t, "api/transfer", "POST", `{"sender":"c", "receiver":"a", "amount": 2}`)
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/accounts", "GET", "")
	assert.NoError(t, err, "GET: api/accounts")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	a = []api.Account{}
	err = json.NewDecoder(res.Body).Decode(&a)
	assert.NoError(t, err, "account list unmarshal")
	assert.Equal(t, a, []api.Account{{"a", 6}, {"b", -6}, {"c", 0}}, "a")

	res, err = testHTTP(t, "api/transfer", "POST", `{"sender":"b", "receiver":"a", "amount": 6}`)
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/accounts", "GET", "")
	assert.NoError(t, err, "GET: api/accounts")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	a = []api.Account{}
	err = json.NewDecoder(res.Body).Decode(&a)
	assert.NoError(t, err, "account list unmarshal")
	assert.Equal(t, a, []api.Account{{"a", 0}, {"b", 0}, {"c", 0}}, "a")

	res, err = testHTTP(t, "api/reset", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")
}
