package cache

import (
	"github.com/allegro/bigcache"
)

// Local struct stores all results in the memory.
type Local struct {
	cache *bigcache.BigCache
}

// Init creates a new BigCache.
func (l *Local) init() error {
	config := bigcache.DefaultConfig(0)
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		return err
	}
	l.cache = cache

	return nil
}

// Set adds a new entry to the cache or overwrite it if already exists.
func (l *Local) Set(key string, data []byte) error {
	err := l.cache.Set(key, data)

	return err
}

// Get returns the value on the key from the cache.
func (l *Local) Get(key string) ([]byte, error) {
	data, err := l.cache.Get(key)
	if err != nil {
		// If not found return our common error.
		if err == bigcache.ErrEntryNotFound {
			err = ErrNotFound
		}
		return nil, err
	}

	return data, nil
}

// Remove deletes an entry from the cache.
func (l *Local) Remove(key string) error {
	err := l.cache.Delete(key)
	if err == bigcache.ErrEntryNotFound {
		err = ErrNotFound
	}

	return err
}

// Clear removes all entries from the cache.
func (l *Local) Clear() error {
	err := l.cache.Reset()

	return err
}
