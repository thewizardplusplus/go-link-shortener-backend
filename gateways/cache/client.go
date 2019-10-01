package cache

import (
	"github.com/go-redis/redis"
)

// Client ...
type Client struct {
	innerClient *redis.Client
}

// NewClient ...
func NewClient(address string) Client {
	innerClient := redis.NewClient(&redis.Options{Addr: address})
	return Client{innerClient: innerClient}
}
