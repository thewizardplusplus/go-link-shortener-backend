// +build integration

package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/caarlos0/env"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

func TestLinkRedirect(test *testing.T) {
	// nolint: lll
	type options struct {
		ServerAddress  string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
		CacheAddress   string `env:"CACHE_ADDRESS" envDefault:"localhost:6379"`
		StorageAddress string `env:"STORAGE_ADDRESS" envDefault:"mongodb://localhost:27017"`
	}

	var opts options
	err := env.Parse(&opts)
	require.NoError(test, err)

	cache := redis.NewClient(&redis.Options{Addr: opts.CacheAddress})
	storage, err := mongo.Connect(
		context.Background(),
		mongooptions.Client().ApplyURI(opts.StorageAddress),
	)
	require.NoError(test, err)

	for _, data := range []struct {
		name         string
		prepare      func(test *testing.T)
		request      *http.Request
		wantStatus   int
		wantLocation string
	}{
		{
			name: "with a link",
			prepare: func(test *testing.T) {
				err := cache.FlushDB().Err()
				require.NoError(test, err)

				_, err = storage.
					Database("go-link-shortener").
					Collection("links").
					DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)

				_, err = storage.
					Database("go-link-shortener").
					Collection("links").
					InsertOne(
						context.Background(),
						entities.Link{Code: "code", URL: "http://example.com/"},
					)
				require.NoError(test, err)
			},
			request: func() *http.Request {
				request, _ := http.NewRequest(
					http.MethodGet,
					opts.ServerAddress+"/redirect/code",
					nil,
				)
				return request
			}(),
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "http://example.com/",
		},
		{
			name: "without a link",
			prepare: func(test *testing.T) {
				err := cache.FlushDB().Err()
				require.NoError(test, err)

				_, err = storage.
					Database("go-link-shortener").
					Collection("links").
					DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)
			},
			request: func() *http.Request {
				request, _ := http.NewRequest(
					http.MethodGet,
					opts.ServerAddress+"/redirect/code",
					nil,
				)
				return request
			}(),
			wantStatus:   http.StatusFound,
			wantLocation: "/error",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test)

			httpClient := &http.Client{
				// disable redirects
				CheckRedirect: func(*http.Request, []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			response, err := httpClient.Do(data.request)
			require.NoError(test, err)

			assert.Equal(test, data.wantStatus, response.StatusCode)
			assert.Equal(test, data.wantLocation, response.Header.Get("Location"))
			assert.Equal(
				test,
				"text/html; charset=utf-8",
				response.Header.Get("Content-Type"),
			)
		})
	}
}
