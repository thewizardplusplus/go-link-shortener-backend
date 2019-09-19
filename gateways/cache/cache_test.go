package cache

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestNewCache(test *testing.T) {
	key := func(link entities.Link) string { return link.Code }
	cache := NewCache("localhost:6379", time.Second, key)

	require.NotNil(test, cache.client)
	assert.Equal(test, "localhost:6379", cache.client.Options().Addr)
	assert.Equal(test, time.Second, cache.expiration)
	assert.Equal(test, getPointer(key), getPointer(cache.key))
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}
