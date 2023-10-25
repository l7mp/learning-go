//go:build localconstructor

package main

import (
	"testing"

	"splitdim/pkg/api"
	"splitdim/pkg/db/local"

	"github.com/stretchr/testify/assert"
)

// TestLocalDBConstructor will create an empty local DB.
func TestLocalDBConstructor(t *testing.T) {
	// this will fail at compile time if something is wrong
	var db api.DataLayer
	db = local.NewDataLayer()
	assert.NotNil(t, db, "not nil")
}
