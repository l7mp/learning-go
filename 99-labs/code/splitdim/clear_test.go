//go:build clear

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"splitdim/pkg/api"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Sprintf("dummy")

// Test through the HTTP layer. Some semi-random transfers and the correponnding account lists.
func TestClear(t *testing.T) {
	// reset!
	res, err := testHTTP(t, "api/reset", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/transfer", "POST", `{"sender":"a", "receiver":"b", "amount": 4}`)
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/clear", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	ts := []api.Transfer{}
	err = json.NewDecoder(res.Body).Decode(&ts)
	assert.NoError(t, err, "clear unmarshal")
	assert.Equal(t, ts, []api.Transfer{{"b", "a", 4}}, "clear transfers")

	// check if accounts are left intact!
	res, err = testHTTP(t, "api/accounts", "GET", "")
	assert.NoError(t, err, "GET: api/accounts")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	a := []api.Account{}
	err = json.NewDecoder(res.Body).Decode(&a)
	assert.NoError(t, err, "account list unmarshal")
	// assert.Equal(t, a, []api.Account{{"a", 4}, {"b", -4}}, "a")
	assert.Len(t, a, 2)
	assert.Contains(t, a, api.Account{"a", 4})
	assert.Contains(t, a, api.Account{"b", -4})

	res, err = testHTTP(t, "api/transfer", "POST", `{"sender":"a", "receiver":"b", "amount": 4}`)
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/clear", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	ts = []api.Transfer{}
	err = json.NewDecoder(res.Body).Decode(&ts)
	assert.NoError(t, err, "clear unmarshal")
	assert.Equal(t, ts, []api.Transfer{{"b", "a", 8}}, "clear transfers ")

	res, err = testHTTP(t, "api/transfer", "POST", `{"sender":"b", "receiver":"c", "amount": 2}`)
	assert.NoError(t, err, "POST: api/transfer")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	res, err = testHTTP(t, "api/clear", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")

	ts = []api.Transfer{}
	err = json.NewDecoder(res.Body).Decode(&ts)
	assert.NoError(t, err, "clear unmarshal")
	assert.Len(t, ts, 2, "clear transfers len")
	assert.Contains(t, ts, api.Transfer{"b", "a", 6}, "clear transfers 1")
	assert.Contains(t, ts, api.Transfer{"c", "a", 2}, "clear transfers 1")

	res, err = testHTTP(t, "api/reset", "GET", "")
	assert.NoError(t, err, "GET: api/rest")
	assert.Equal(t, http.StatusOK, res.StatusCode, "status")
}
