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
