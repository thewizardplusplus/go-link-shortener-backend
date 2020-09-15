// +build integration
// +build bulky

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"testing"

	"github.com/caarlos0/env"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/storage"
	"go.etcd.io/etcd/clientv3"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLinkCreating_bulky(test *testing.T) {
	// nolint: lll
	type options struct {
		ServerAddress  string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
		CacheAddress   string `env:"CACHE_ADDRESS" envDefault:"localhost:6379"`
		StorageAddress string `env:"STORAGE_ADDRESS" envDefault:"mongodb://localhost:27017"`
		Counter        struct {
			Address string `env:"COUNTER_ADDRESS" envDefault:"localhost:2379"`
			Count   int    `env:"COUNTER_COUNT" envDefault:"2"`
			Chunk   uint64 `env:"COUNTER_CHUNK" envDefault:"1000"`
			Range   uint64 `env:"COUNTER_RANGE" envDefault:"1000000000"`
		}
	}

	var opts options
	err := env.Parse(&opts)
	require.NoError(test, err)

	cache := redis.NewClient(&redis.Options{Addr: opts.CacheAddress})
	storage, err :=
		storage.NewClient(opts.StorageAddress, "go-link-shortener", "links")
	require.NoError(test, err)

	counter, err := clientv3.New(clientv3.Config{
		Endpoints: []string{opts.Counter.Address},
	})
	require.NoError(test, err)

	for _, data := range []struct {
		name    string
		prepare func(test *testing.T) (preparedData interface{})
		count   int
		check   func(test *testing.T, preparedData interface{}, codes []uint64)
	}{
		{
			name: "success",
			prepare: func(test *testing.T) (preparedData interface{}) {
				err := cache.FlushDB().Err()
				require.NoError(test, err)

				_, err = storage.Collection().DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)

				var counters []uint64
				for i := 0; i < opts.Counter.Count; i++ {
					name := fmt.Sprintf("distributed-counter-%d", i)

					_, err := counter.Put(context.Background(), name, "")
					require.NoError(test, err)

					response, err := counter.Get(context.Background(), name)
					require.NoError(test, err)
					require.NotNil(test, response.Kvs)

					counters = append(counters, uint64(response.Kvs[0].Version))
				}

				return counters
			},
			count: 10000,
			check: func(test *testing.T, preparedData interface{}, codes []uint64) {
				var testedCodeCount int
				for i := 0; i < opts.Counter.Count; i++ {
					var counterCodes []uint64
					minimum := uint64(i) * opts.Counter.Range
					maximum := uint64(i+1) * opts.Counter.Range
					for _, code := range codes {
						if code >= minimum && code < maximum {
							counterCodes = append(counterCodes, code)
						}
					}
					if len(counterCodes) == 0 {
						continue
					}

					// check that the row start differs from the counter
					// by no more than the chunk size
					assert.InDelta(
						test,
						preparedData.([]uint64)[i]*opts.Counter.Chunk+minimum,
						counterCodes[0],
						float64(opts.Counter.Chunk),
					)

					// check that codes in the row increase and differ
					// by either 1 or the chunk size
					for j := 1; j < len(counterCodes); j++ {
						delta := counterCodes[j] - counterCodes[j-1]
						assert.Contains(test, []uint64{1, opts.Counter.Chunk + 1}, delta)
					}

					testedCodeCount += len(counterCodes)
				}

				// check that all codes are in the range of any counter
				assert.Equal(test, len(codes), testedCodeCount)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			preparedData := data.prepare(test)

			var codes []uint64
			for i := 0; i < data.count; i++ {
				// print progress on its every one percent
				if i%(data.count/100) == 0 {
					log.Printf("%d/%d links were created", i, data.count)
				}

				request, err := http.NewRequest(
					http.MethodPost,
					opts.ServerAddress+"/api/v1/links/",
					bytes.NewBufferString(fmt.Sprintf(
						`{"URL":"http://example.com/pages/%d"}`,
						i,
					)),
				)
				require.NoError(test, err)

				response, err := http.DefaultClient.Do(request)
				require.NoError(test, err)
				defer response.Body.Close()

				var link entities.Link
				err = json.NewDecoder(response.Body).Decode(&link)
				require.NoError(test, err)

				code := new(big.Int)
				_, ok := code.SetString(link.Code, 62)
				require.True(test, ok)

				codes = append(codes, code.Uint64())
			}

			data.check(test, preparedData, codes)
		})
	}
}
