package cache

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// KeyExtractor ...
type KeyExtractor func(link entities.Link) string

// Cache ...
type Cache struct {
	client     *redis.Client
	expiration time.Duration
	key        KeyExtractor
}

// NewCache ...
func NewCache(
	address string,
	expiration time.Duration,
	key KeyExtractor,
) Cache {
	client := redis.NewClient(&redis.Options{Addr: address})
	return Cache{client, expiration, key}
}
