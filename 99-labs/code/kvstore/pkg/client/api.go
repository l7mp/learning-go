package client

import "kvstore/pkg/api"

// Client is a generic client interface to the key-value store.
type Client interface {
	// Get returns the the value and version stored for the given key, or an error if something goes wrong.
	Get(key string) (api.VersionedValue, error)
	// Put tries to insert the given key-value pair with the specified version into the store.
	Put(vkv api.VersionedKeyValue) error
	// List returns all values stored in the database.
	List() ([]api.VersionedKeyValue, error)
	// Reset removes all key-value pairs.
	Reset() error
}
