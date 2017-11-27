package storage

import "github.com/erans/upupaway/config"

var activeStorage Storage

// Storage defines the interface for the simplest get/set key operations
type Storage interface {
	Get(key string) string
	Set(key string, value string)
}

// InitActiveStorage sets up the active storage engine
func InitActiveStorage(cfg *config.Config) {
	if activeStorage == nil {
		switch cfg.Storage.ActiveStorage {
		case "redis":
			activeStorage = NewRedisStorage(cfg)
		case "inmemory":
			activeStorage = NewInMemoryStorage(cfg)
		case "memcache":
			activeStorage = NewMemCacheStorage(cfg)
		}
	}
}

// GetActiveStorage returns the active storage engine
func GetActiveStorage() Storage {
	return activeStorage
}
