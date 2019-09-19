// +build integration

package cache

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

func TestCache_GetLink(test *testing.T) {
	type fields struct {
		client *redis.Client
	}
	type args struct {
		query string
	}

	for _, data := range []struct {
		name     string
		fields   fields
		prepare  func(test *testing.T, client *redis.Client)
		args     args
		wantLink entities.Link
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test, data.fields.client)

			cache := Cache{
				client: data.fields.client,
			}
			gotLink, gotErr := cache.GetLink(data.args.query)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
