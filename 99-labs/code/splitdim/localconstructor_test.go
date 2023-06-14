//go:build localconstructor

package main

import (
	"testing"

	"splitdim/pkg/api"
	"splitdim/pkg/db/local"

	"github.com/stretchr/testify/assert"
)

// TestAPI will simply create the API structs and JSON encode/decode each.
func TestLocalDBConstrurctor(t *testing.T) {
	// this will fail at compile time if something is wrong
	var db api.DataLayer
	db = local.NewDataLayer()
	assert.NotNil(t, db, "not nil")
}
