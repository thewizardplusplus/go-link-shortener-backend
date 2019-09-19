// +build integration

package cache

import (
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

var (
	address string
)

func init() {
	var ok bool
	if address, ok = os.LookupEnv("REDIS_URL"); !ok {
		address = "localhost:6379"
	}
}

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
		{
			name: "success",
			fields: fields{
				client: redis.NewClient(&redis.Options{Addr: address}),
			},
			prepare: func(test *testing.T, client *redis.Client) {
				err := client.Set("query", `{"Code":"code","URL":"url"}`, 0).Err()
				require.NoError(test, err)
			},
			args:     args{"query"},
			wantLink: entities.Link{Code: "code", URL: "url"},
			wantErr:  assert.NoError,
		},
		{
			name: "error without data",
			fields: fields{
				client: redis.NewClient(&redis.Options{Addr: address}),
			},
			prepare: func(test *testing.T, client *redis.Client) {
				err := client.Del("query").Err()
				require.NoError(test, err)
			},
			args:     args{"query"},
			wantLink: entities.Link{},
			wantErr:  assert.Error,
		},
		{
			name: "error with incorrect data",
			fields: fields{
				client: redis.NewClient(&redis.Options{Addr: address}),
			},
			prepare: func(test *testing.T, client *redis.Client) {
				err := client.Set("query", "incorrect", 0).Err()
				require.NoError(test, err)
			},
			args:     args{"query"},
			wantLink: entities.Link{},
			wantErr:  assert.Error,
		},
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

func TestCache_SetLink(test *testing.T) {
	type fields struct {
		client     *redis.Client
		expiration time.Duration
		key        KeyExtractor
	}
	type args struct {
		link entities.Link
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, client *redis.Client)
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			cache := Cache{
				client:     data.fields.client,
				expiration: data.fields.expiration,
				key:        data.fields.key,
			}
			gotErr := cache.SetLink(data.args.link)

			data.wantErr(test, gotErr)
			data.check(test, data.fields.client)
		})
	}
}
