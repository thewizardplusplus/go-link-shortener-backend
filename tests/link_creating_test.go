// +build integration

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	"go.etcd.io/etcd/clientv3"
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
		Counter        struct {
			Address string `env:"COUNTER_ADDRESS" envDefault:"localhost:2379"`
			Count   int    `env:"COUNTER_COUNT" envDefault:"2"`
		}
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

	counter, err := clientv3.New(clientv3.Config{
		Endpoints: []string{opts.Counter.Address},
	})
	require.NoError(test, err)

	for _, data := range []struct {
		name            string
		prepare         func(test *testing.T) (preparedData interface{})
		request         *http.Request
		wantStatus      int
		wantCodePattern *regexp.Regexp
		wantURL         string
		check           func(
			test *testing.T,
			preparedData interface{},
			expectedLink entities.Link,
		)
	}{
		{
			name: "without a link inside the cache and the storage",
			prepare: func(test *testing.T) (preparedData interface{}) {
				err := cache.FlushDB().Err()
				require.NoError(test, err)

				_, err = storage.
					Database("go-link-shortener").
					Collection("links").
					DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)

				var counters []uint64
				for i := 0; i < opts.Counter.Count; i++ {
					name := fmt.Sprintf("distributed-counter-%d", i)
					response, err := counter.Get(context.Background(), name)
					require.NoError(test, err)
					require.NotNil(test, response.Kvs)

					counters = append(counters, uint64(response.Kvs[0].Version))
				}

				return counters
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
			wantCodePattern: regexp.MustCompile(`[\da-zA-Z]+`), // base 62
			wantURL:         "http://example.com/",
			check: func(
				test *testing.T,
				preparedData interface{},
				expectedLink entities.Link,
			) {
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

				var counters []uint64
				for i := 0; i < opts.Counter.Count; i++ {
					name := fmt.Sprintf("distributed-counter-%d", i)
					response, err := counter.Get(context.Background(), name)
					require.NoError(test, err)
					require.NotNil(test, response.Kvs)

					counters = append(counters, uint64(response.Kvs[0].Version))
				}

				preparedDataVariants := [][]uint64{preparedData.([]uint64)}
				for i := 0; i < opts.Counter.Count; i++ {
					preparedDataCopy := make([]uint64, opts.Counter.Count)
					copy(preparedDataCopy, preparedData.([]uint64))

					preparedDataCopy[i]++

					preparedDataVariants = append(preparedDataVariants, preparedDataCopy)
				}

				assert.Equal(test, expectedLink, link)
				assert.Contains(test, preparedDataVariants, counters)
			},
		},
		{
			name: "with a link inside the cache",
			prepare: func(test *testing.T) (preparedData interface{}) {
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

				var counters []uint64
				for i := 0; i < opts.Counter.Count; i++ {
					name := fmt.Sprintf("distributed-counter-%d", i)
					response, err := counter.Get(context.Background(), name)
					require.NoError(test, err)
					require.NotNil(test, response.Kvs)

					counters = append(counters, uint64(response.Kvs[0].Version))
				}

				return counters
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
			check: func(
				test *testing.T,
				preparedData interface{},
				expectedLink entities.Link,
			) {
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

				var counters []uint64
				for i := 0; i < opts.Counter.Count; i++ {
					name := fmt.Sprintf("distributed-counter-%d", i)
					response, err := counter.Get(context.Background(), name)
					require.NoError(test, err)
					require.NotNil(test, response.Kvs)

					counters = append(counters, uint64(response.Kvs[0].Version))
				}

				assert.Equal(test, expectedLink, link)
				assert.InDelta(test, time.Hour, duration, float64(10*time.Second))
				assert.Equal(test, preparedData, counters)
			},
		},
		{
			name: "with a link inside the storage",
			prepare: func(test *testing.T) (preparedData interface{}) {
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

				var counters []uint64
				for i := 0; i < opts.Counter.Count; i++ {
					name := fmt.Sprintf("distributed-counter-%d", i)
					response, err := counter.Get(context.Background(), name)
					require.NoError(test, err)
					require.NotNil(test, response.Kvs)

					counters = append(counters, uint64(response.Kvs[0].Version))
				}

				return counters
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
			check: func(
				test *testing.T,
				preparedData interface{},
				expectedLink entities.Link,
			) {
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

				var counters []uint64
				for i := 0; i < opts.Counter.Count; i++ {
					name := fmt.Sprintf("distributed-counter-%d", i)
					response, err := counter.Get(context.Background(), name)
					require.NoError(test, err)
					require.NotNil(test, response.Kvs)

					counters = append(counters, uint64(response.Kvs[0].Version))
				}

				assert.Equal(test, expectedLink, link)
				assert.Equal(test, preparedData, counters)
			},
		},
		{
			name: "with a link inside the cache and the storage",
			prepare: func(test *testing.T) (preparedData interface{}) {
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

				var counters []uint64
				for i := 0; i < opts.Counter.Count; i++ {
					name := fmt.Sprintf("distributed-counter-%d", i)
					response, err := counter.Get(context.Background(), name)
					require.NoError(test, err)
					require.NotNil(test, response.Kvs)

					counters = append(counters, uint64(response.Kvs[0].Version))
				}

				return counters
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
			check: func(
				test *testing.T,
				preparedData interface{},
				expectedLink entities.Link,
			) {
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

				var counters []uint64
				for i := 0; i < opts.Counter.Count; i++ {
					name := fmt.Sprintf("distributed-counter-%d", i)
					response, err := counter.Get(context.Background(), name)
					require.NoError(test, err)
					require.NotNil(test, response.Kvs)

					counters = append(counters, uint64(response.Kvs[0].Version))
				}

				assert.Equal(test, expectedLink, link)
				assert.InDelta(test, time.Hour, duration, float64(10*time.Second))
				assert.Equal(test, preparedData, counters)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			preparedData := data.prepare(test)

			response, err := http.DefaultClient.Do(data.request)
			require.NoError(test, err)
			defer response.Body.Close()

			var link entities.Link
			json.NewDecoder(response.Body).Decode(&link)

			assert.Equal(test, data.wantStatus, response.StatusCode)
			assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
			assert.Regexp(test, data.wantCodePattern, link.Code)
			assert.Equal(test, data.wantURL, link.URL)
			data.check(test, preparedData, link)
		})
	}
}
