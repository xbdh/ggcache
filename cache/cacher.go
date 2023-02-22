package cache

import "time"

type Cacher interface {
	// Get returns the value for the given key.
	Get(key []byte) ([]byte, error)
	// sets the value for the given key.
	Set(key []byte, value []byte, ttl time.Duration) error
	// Delete removes the value for the given key.
	Delete(key []byte) error
	// has returns true if the given key exists.
	Has(key []byte) bool
}
