package client

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"kvstore/pkg/api"
	"kvstore/pkg/server"
)

func TestClient(t *testing.T) {
	f, err := os.CreateTemp("", "sample")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	s, err := server.NewServer(f.Name())
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addr := ":8081"
	assert.NoError(t, s.Run(ctx, addr))

	c := NewClient(addr)

	// empty
	err = c.Reset()
	assert.NoError(t, err)

	// indeed?
	vv, err := c.Get("dummy")
	assert.NoError(t, err)
	assert.Equal(t, "", vv.Value)
	assert.Equal(t, 0, vv.Version)

	l, err := c.List()
	assert.NoError(t, err)
	assert.Len(t, l, 0)

	// add something
	err = c.Put(api.VersionedKeyValue{Key: "key1", VersionedValue: api.VersionedValue{Value: "a", Version: 0}})
	assert.NoError(t, err)

	vv, err = c.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "a", vv.Value)
	assert.Equal(t, 1, vv.Version)

	l, err = c.List()
	assert.NoError(t, err)
	assert.Len(t, l, 1)
	assert.Equal(t, "key1", l[0].Key)
	assert.Equal(t, "a", l[0].Value)
	assert.Equal(t, 1, l[0].Version)

	// wrong version fails
	err = c.Put(api.VersionedKeyValue{Key: "key1", VersionedValue: api.VersionedValue{Value: "b", Version: 0}})
	assert.Error(t, err)

	// update works
	err = c.Put(api.VersionedKeyValue{Key: "key1", VersionedValue: api.VersionedValue{Value: "b", Version: 1}})
	assert.NoError(t, err)

	vv, err = c.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "b", vv.Value)
	assert.Equal(t, 2, vv.Version)

	l, err = c.List()
	assert.NoError(t, err)
	assert.Len(t, l, 1)
	assert.Equal(t, "key1", l[0].Key)
	assert.Equal(t, "b", l[0].Value)
	assert.Equal(t, 2, l[0].Version)

	// add something else
	err = c.Put(api.VersionedKeyValue{Key: "key2", VersionedValue: api.VersionedValue{Value: "c", Version: 0}})
	assert.NoError(t, err)

	vv, err = c.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "b", vv.Value)
	assert.Equal(t, 2, vv.Version)

	vv, err = c.Get("key2")
	assert.NoError(t, err)
	assert.Equal(t, "c", vv.Value)
	assert.Equal(t, 1, vv.Version)

	l, err = c.List()
	assert.NoError(t, err)
	assert.Len(t, l, 2)

	// empty
	err = c.Reset()
	assert.NoError(t, err)

	// indeed?
	vv, err = c.Get("dummy")
	assert.NoError(t, err)
	assert.Equal(t, "", vv.Value)
	assert.Equal(t, 0, vv.Version)

	l, err = c.List()
	assert.NoError(t, err)
	assert.Len(t, l, 0)
}
