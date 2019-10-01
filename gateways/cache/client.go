package cache

import (
	"github.com/go-redis/redis"
)

// Client ...
type Client struct {
	innerClient *redis.Client
}
