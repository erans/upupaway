package storage

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/erans/upupaway/config"
)

// MemCacheStorage is the implementation of a memcache based storage system
type MemCacheStorage struct {
	Cache             *memcache.Client
	DefaultExpiration int32
}

// Get returns the value of the specified key
func (c *MemCacheStorage) Get(key string) string {
	if item, err := c.Cache.Get(key); err == nil {
		return string(item.Value[:])
	}

	return ""
}

// Set allows setting a value for the specified key
func (c *MemCacheStorage) Set(key string, value string) {
	c.Cache.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: c.DefaultExpiration})
}

// NewMemCacheStorage returns a new memcache storage object
func NewMemCacheStorage(cfg *config.Config) *MemCacheStorage {
	return &MemCacheStorage{
		Cache:             memcache.New(fmt.Sprintf("%s:%d", cfg.Storage.MemCache.Host, cfg.Storage.MemCache.Port)),
		DefaultExpiration: cfg.Storage.MemCache.Expiration,
	}
}
