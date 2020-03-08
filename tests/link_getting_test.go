// +build integration

package tests

// nolint: lll
import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/http/handlers/presenters"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

func TestLinkGetting(test *testing.T) {
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
		name       string
		prepare    func(test *testing.T)
		request    *http.Request
		response   interface{}
		wantStatus int
		wantCode   string
		wantURL    string
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
					http.MethodGet,
					opts.ServerAddress+"/api/v1/links/code",
					nil,
				)
				return request
			}(),
			response:   new(presenters.ErrorResponse),
			wantStatus: http.StatusNotFound,
		},
		{
			name: "with a link inside the cache",
			prepare: func(test *testing.T) {
				err := cache.FlushDB().Err()
				require.NoError(test, err)

				err = cache.
					Set(
						"code",
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
					http.MethodGet,
					opts.ServerAddress+"/api/v1/links/code",
					nil,
				)
				return request
			}(),
			response:   new(entities.Link),
			wantStatus: http.StatusOK,
			wantCode:   "code",
			wantURL:    "http://example.com/",
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
					http.MethodGet,
					opts.ServerAddress+"/api/v1/links/code",
					nil,
				)
				return request
			}(),
			response:   new(entities.Link),
			wantStatus: http.StatusOK,
			wantCode:   "code",
			wantURL:    "http://example.com/",
		},
		{
			name: "with a link inside the cache and the storage",
			prepare: func(test *testing.T) {
				err := cache.FlushDB().Err()
				require.NoError(test, err)

				err = cache.
					Set(
						"code",
						`{"Code":"code","URL":"http://example.com/1"}`,
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
						entities.Link{Code: "code", URL: "http://example.com/2"},
					)
				require.NoError(test, err)
			},
			request: func() *http.Request {
				request, _ := http.NewRequest(
					http.MethodGet,
					opts.ServerAddress+"/api/v1/links/code",
					nil,
				)
				return request
			}(),
			response:   new(entities.Link),
			wantStatus: http.StatusOK,
			wantCode:   "code",
			wantURL:    "http://example.com/1",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test)

			response, err := http.DefaultClient.Do(data.request)
			require.NoError(test, err)
			defer response.Body.Close()

			json.NewDecoder(response.Body).Decode(data.response)

			assert.Equal(test, data.wantStatus, response.StatusCode)
			assert.Equal(test, "application/json", response.Header.Get("Content-Type"))
			switch response := data.response.(type) {
			case *entities.Link:
				assert.Equal(test, data.wantCode, response.Code)
				assert.Equal(test, data.wantURL, response.URL)
			case *presenters.ErrorResponse:
				assert.NotEmpty(test, response.Error)
			default:
				test.Errorf("incorrect response type: %T", response)
			}
		})
	}
}
