// +build integration

package cache

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

var (
	address string
)

// nolint: gochecknoinits
func init() {
	var ok bool
	if address, ok = os.LookupEnv("REDIS_URL"); !ok {
		address = "localhost:6379"
	}
}

func TestLinkGetter_GetLink(test *testing.T) {
	type fields struct {
		Client Client
	}
	type args struct {
		query string
	}

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
				Client: NewClient(address),
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
				Client: NewClient(address),
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
				Client: NewClient(address),
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
