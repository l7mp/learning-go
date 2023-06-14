//go:build api

package main

import (
	"encoding/json"
	"testing"

	"splitdim/pkg/api"

	"github.com/stretchr/testify/assert"
)

type dummyDB struct{}

func (d *dummyDB) Transfer(t api.Transfer) error       { return nil }
func (d *dummyDB) AccountList() ([]api.Account, error) { return []api.Account{}, nil }
func (d *dummyDB) Clear() ([]api.Transfer, error)      { return []api.Transfer{}, nil }
func (d *dummyDB) Reset() error                        { return nil }

// TestAPI will simply create the API structs and JSON encode/decode each.
func TestAPI(t *testing.T) {
	// this will fail at compile time if something is wrong
	tr := api.Transfer{
		Sender:   "a",
		Receiver: "b",
		Amount:   1,
	}

	// unmarshal test to check JSON tags
	j := []byte(`{"sender":"a","receiver":"b", "amount": 1}`)

	tr2 := api.Transfer{}
	err := json.Unmarshal(j, &tr2)
	assert.NoError(t, err, "transfer JSON error")
	assert.Equal(t, tr, tr2, "transfer unmarshal equal")

	// this will fail at compile time if something is wrong
	a := api.Account{
		Holder:  "a",
		Balance: 1,
	}

	// unmarshal test to check JSON tags
	j = []byte(`{"holder":"a", "balance": 1}`)

	a2 := api.Account{}
	err = json.Unmarshal(j, &a2)
	assert.NoError(t, err, "account JSON error")
	assert.Equal(t, a, a2, "transfer unmarshal equal")

	// API test: fail at compile time
	d := &dummyDB{}
	var x api.DataLayer
	x = d
	assert.NotNil(t, x, "noerror")
}
