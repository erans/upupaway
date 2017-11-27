package storage

import (
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/erans/upupaway/config"
)

// InMemoryStorage is the implementation of an in-memory storage system - mostly good for debug/local usage
type InMemoryStorage struct {
	Cache *cache.Cache
}

// Get returns the value of the specified key
func (c *InMemoryStorage) Get(key string) string {
	if v, ok := c.Cache.Get(key); ok {
		return v.(string)
	}

	return ""
}

// Set allows setting a value for the specified key
func (c *InMemoryStorage) Set(key string, value string) {
	c.Cache.SetDefault(key, value)
}

// NewInMemoryStorage returns a new in-memory storage object
func NewInMemoryStorage(cfg *config.Config) *InMemoryStorage {
	c := cache.New(time.Duration(cfg.Storage.InMemory.Expiration)*time.Second, 10*time.Minute)
	return &InMemoryStorage{
		Cache: c,
	}
}
