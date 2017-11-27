package storage

import (
	"fmt"
	"time"

	"github.com/erans/upupaway/config"
	"github.com/go-redis/redis"
)

// RedisStorage is the implementation of a redis based storage system
type RedisStorage struct {
	Client            *redis.Client
	DefaultExpiration int
}

// Get returns the value of the specified key
func (r *RedisStorage) Get(key string) string {
	result, _ := r.Client.Get(key).Result()
	return result
}

// Set allows setting a value for the specified key
func (r *RedisStorage) Set(key string, value string) {
	r.Client.Set(key, value, time.Second*time.Duration(r.DefaultExpiration))
}

// NewRedisStorage returns a new redis storage object
func NewRedisStorage(cfg *config.Config) *RedisStorage {
	options := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Storage.Redis.Host, cfg.Storage.Redis.Port),
		DB:   cfg.Storage.Redis.DB,
	}

	client := redis.NewClient(options)

	return &RedisStorage{
		Client:            client,
		DefaultExpiration: cfg.Storage.Redis.Expiration,
	}
}
