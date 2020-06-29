package cache

import (
	"errors"
	"fmt"
)

// ErrNotFound means the key is not in the cache.
var ErrNotFound = errors.New("entry not found")

// Cache is an interface for interacting with various cache backend.
type Cache interface {
	init() error                       // Initalizes the backend.
	Set(key string, data []byte) error // Set a key-value to the cache.
	Get(key string) ([]byte, error)    // Should return `ErrNotFound` if the key not in the cache.
	Remove(key string) error           // Should return `ErrNotFound` if the key not in the cache.
	Clear() error                      // Clear the whole cache.
}

// NewCache returns a newly created `local` or `redis` cache system.
func NewCache(cacheType string) (Cache, error) {
	switch cacheType {
	case "local":
		local := &Local{}
		if err := local.init(); err != nil {
			return nil, err
		}
		return local, nil

	default:
		return nil, fmt.Errorf(`cache "%s" is not implemented`, cacheType)
	}
}
