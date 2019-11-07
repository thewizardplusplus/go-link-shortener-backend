// +build integration

package counter

import (
	"context"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounter_NextCountChunk(test *testing.T) {
	type options struct {
		CounterAddress string `env:"COUNTER_ADDRESS" envDefault:"localhost:2379"`
	}
	type fields struct {
		makeClient func(test *testing.T) Client
		name       string
	}

	var opts options
	err := env.Parse(&opts)
	require.NoError(test, err)

	for _, data := range []struct {
		name    string
		prepare func(test *testing.T, counter Counter) (preparedData interface{})
		fields  fields
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, preparedData interface{}, gotChunk uint64)
	}{
		{
			name: "success",
			prepare: func(test *testing.T, counter Counter) (preparedData interface{}) {
				context := context.Background()
				response, err := counter.Client.innerClient.Get(context, "counter")
				require.NoError(test, err)

				return uint64(response.Header.Revision)
			},
			fields: fields{
				makeClient: func(test *testing.T) Client {
					client, err := NewClient(opts.CounterAddress)
					require.NoError(test, err)

					return client
				},
				name: "counter",
			},
			check: func(test *testing.T, preparedData interface{}, gotChunk uint64) {
				assert.Equal(test, preparedData.(uint64)+1, gotChunk)
			},
			wantErr: assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := data.fields.makeClient(test)
			counter := Counter{
				Client: client,
				Name:   data.fields.name,
			}
			preparedData := data.prepare(test, counter)
			gotChunk, gotErr := counter.NextCountChunk()

			data.wantErr(test, gotErr)
			data.check(test, preparedData, gotChunk)
		})
	}
}
