// +build integration

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

func TestLinkCreating(test *testing.T) {
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
		name            string
		prepare         func(test *testing.T)
		request         *http.Request
		wantStatus      int
		wantCodePattern *regexp.Regexp
		wantURL         string
		check           func(test *testing.T, expectedLink entities.Link)
	}{
		{
			name: "without a link inside the cache and the storage",
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
					http.MethodPost,
					opts.ServerAddress+"/api/v1/links/",
					bytes.NewBufferString(`{"URL":"http://example.com/"}`),
				)
				return request
			}(),
			wantStatus: http.StatusOK,
			wantCodePattern: regexp.MustCompile(
				`(?i)^[\da-f]{8}(-[\da-f]{4}){3}-[\da-f]{12}$`,
			),
			wantURL: "http://example.com/",
			check: func(test *testing.T, expectedLink entities.Link) {
				for _, key := range []string{expectedLink.Code, expectedLink.URL} {
					data, err := cache.Get(key).Result()
					require.NoError(test, err)

					duration, err := cache.TTL(key).Result()
					require.NoError(test, err)

					var link entities.Link
					json.NewDecoder(strings.NewReader(data)).Decode(&link)

					assert.Equal(test, expectedLink, link)
					assert.InDelta(test, time.Hour, duration, float64(10*time.Second))
				}

				var link entities.Link
				err := storage.
					Database("go-link-shortener").
					Collection("links").
					FindOne(context.Background(), bson.M{"code": expectedLink.Code}).
					Decode(&link)
				require.NoError(test, err)

				assert.Equal(test, expectedLink, link)
			},
		},
		{
			name: "with a link inside the cache",
			prepare: func(test *testing.T) {
				err := cache.FlushDB().Err()
				require.NoError(test, err)

				err = cache.
					Set(
						"http://example.com/",
						`{"Code":"code","URL":"http://example.com/"}`,
						time.Hour,
					).
					Err()
				require.NoError(test, err)

				_, err = storage.
					Database("go-link-shortener").
					Collection("links").
					DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)
			},
			request: func() *http.Request {
				request, _ := http.NewRequest(
					http.MethodPost,
					opts.ServerAddress+"/api/v1/links/",
					bytes.NewBufferString(`{"URL":"http://example.com/"}`),
				)
				return request
			}(),
			wantStatus:      http.StatusOK,
			wantCodePattern: regexp.MustCompile("code"),
			wantURL:         "http://example.com/",
			check: func(test *testing.T, expectedLink entities.Link) {
				_, err := cache.Get(expectedLink.Code).Result()
				require.Equal(test, redis.Nil, err)

				data, err := cache.Get(expectedLink.URL).Result()
				require.NoError(test, err)

				duration, err := cache.TTL(expectedLink.URL).Result()
				require.NoError(test, err)

				var link entities.Link
				json.NewDecoder(strings.NewReader(data)).Decode(&link)

				err = storage.
					Database("go-link-shortener").
					Collection("links").
					FindOne(context.Background(), bson.M{"code": expectedLink.Code}).
					Err()
				require.Equal(test, mongo.ErrNoDocuments, err)

				assert.Equal(test, expectedLink, link)
				assert.InDelta(test, time.Hour, duration, float64(10*time.Second))
			},
		},
		{
			name: "with a link inside the storage",
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
					http.MethodPost,
					opts.ServerAddress+"/api/v1/links/",
					bytes.NewBufferString(`{"URL":"http://example.com/"}`),
				)
				return request
			}(),
			wantStatus:      http.StatusOK,
			wantCodePattern: regexp.MustCompile("code"),
			wantURL:         "http://example.com/",
			check: func(test *testing.T, expectedLink entities.Link) {
				_, err := cache.Get(expectedLink.Code).Result()
				assert.Equal(test, redis.Nil, err)

				_, err = cache.Get(expectedLink.URL).Result()
				assert.Equal(test, redis.Nil, err)

				var link entities.Link
				err = storage.
					Database("go-link-shortener").
					Collection("links").
					FindOne(context.Background(), bson.M{"code": expectedLink.Code}).
					Decode(&link)
				require.NoError(test, err)

				assert.Equal(test, expectedLink, link)
			},
		},
		{
			name: "with a link inside the cache and the storage",
			prepare: func(test *testing.T) {
				err := cache.FlushDB().Err()
				require.NoError(test, err)

				err = cache.
					Set(
						"http://example.com/",
						`{"Code":"code #1","URL":"http://example.com/"}`,
						time.Hour,
					).
					Err()
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
						entities.Link{Code: "code #2", URL: "http://example.com/"},
					)
				require.NoError(test, err)
			},
			request: func() *http.Request {
				request, _ := http.NewRequest(
					http.MethodPost,
					opts.ServerAddress+"/api/v1/links/",
					bytes.NewBufferString(`{"URL":"http://example.com/"}`),
				)
				return request
			}(),
			wantStatus:      http.StatusOK,
			wantCodePattern: regexp.MustCompile("code #1"),
			wantURL:         "http://example.com/",
			check: func(test *testing.T, expectedLink entities.Link) {
				_, err := cache.Get(expectedLink.Code).Result()
				assert.Equal(test, redis.Nil, err)

				data, err := cache.Get(expectedLink.URL).Result()
				require.NoError(test, err)

				duration, err := cache.TTL(expectedLink.URL).Result()
				require.NoError(test, err)

				var link entities.Link
				json.NewDecoder(strings.NewReader(data)).Decode(&link)

				err = storage.
					Database("go-link-shortener").
					Collection("links").
					FindOne(context.Background(), bson.M{"code": expectedLink.Code}).
					Err()
				assert.Equal(test, mongo.ErrNoDocuments, err)

				assert.Equal(test, expectedLink, link)
				assert.InDelta(test, time.Hour, duration, float64(10*time.Second))
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test)

			response, err := http.DefaultClient.Do(data.request)
			require.NoError(test, err)
			defer response.Body.Close()

			var link entities.Link
			json.NewDecoder(response.Body).Decode(&link)

			assert.Equal(test, data.wantStatus, response.StatusCode)
			assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
			assert.Regexp(test, data.wantCodePattern, link.Code)
			assert.Equal(test, data.wantURL, link.URL)
			data.check(test, link)
		})
	}
}
