package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
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

// GetLink ...
func (cache Cache) GetLink(query string) (entities.Link, error) {
	data, err := cache.client.Get(query).Result()
	if err != nil {
		return entities.Link{}, errors.Wrap(err, "unable to get the link from Redis")
	}

	var link entities.Link
	if err := json.Unmarshal([]byte(data), &link); err != nil {
		return entities.Link{},
			errors.Wrap(err, "unable to unmarshal the link from Redis")
	}

	return link, nil
}
