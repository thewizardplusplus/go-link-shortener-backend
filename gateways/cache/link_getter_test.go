// +build integration

package cache

import (
	"database/sql"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

func TestLinkGetter_GetLink(test *testing.T) {
	type options struct {
		CacheAddress string `env:"CACHE_ADDRESS" envDefault:"localhost:6379"`
	}
	type fields struct {
		Client Client
	}
	type args struct {
		query string
	}

	var opts options
	err := env.Parse(&opts)
	require.NoError(test, err)

	for _, data := range []struct {
		name     string
		fields   fields
		prepare  func(test *testing.T, client Client)
		args     args
		wantLink entities.Link
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				Client: NewClient(opts.CacheAddress),
			},
			prepare: func(test *testing.T, client Client) {
				err := client.innerClient.
					Set("query", `{"Code":"code","URL":"url"}`, 0).
					Err()
				require.NoError(test, err)
			},
			args:     args{"query"},
			wantLink: entities.Link{Code: "code", URL: "url"},
			wantErr:  assert.NoError,
		},
		{
			name: "error without data",
			fields: fields{
				Client: NewClient(opts.CacheAddress),
			},
			prepare: func(test *testing.T, client Client) {
				err := client.innerClient.Del("query").Err()
				require.NoError(test, err)
			},
			args:     args{"query"},
			wantLink: entities.Link{},
			wantErr: func(test assert.TestingT, err error, args ...interface{}) bool {
				return assert.Equal(test, sql.ErrNoRows, err, args)
			},
		},
		{
			name: "error with incorrect data",
			fields: fields{
				Client: NewClient(opts.CacheAddress),
			},
			prepare: func(test *testing.T, client Client) {
				err := client.innerClient.Set("query", "incorrect", 0).Err()
				require.NoError(test, err)
			},
			args:     args{"query"},
			wantLink: entities.Link{},
			wantErr:  assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			data.prepare(test, data.fields.Client)

			cache := LinkGetter{
				Client: data.fields.Client,
			}
			gotLink, gotErr := cache.GetLink(data.args.query)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
