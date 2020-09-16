// +build integration

package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLinkGetter_GetLink(test *testing.T) {
	// nolint: lll
	type options struct {
		StorageAddress string `env:"STORAGE_ADDRESS" envDefault:"mongodb://localhost:27017"`
	}
	type fields struct {
		makeClient func(test *testing.T) Client
		keyField   string
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
		prepare  func(test *testing.T, getter LinkGetter)
		args     args
		wantLink entities.Link
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				makeClient: func(test *testing.T) Client {
					client, err := NewClient(opts.StorageAddress, "database", "collection")
					require.NoError(test, err)

					return client
				},
				keyField: CodeLinkField,
			},
			prepare: func(test *testing.T, getter LinkGetter) {
				_, err := getter.Client.
					Collection().
					DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)

				_, err = getter.Client.
					Collection().
					InsertOne(context.Background(), entities.Link{Code: "code", URL: "url"})
				require.NoError(test, err)
			},
			args:     args{"code"},
			wantLink: entities.Link{Code: "code", URL: "url"},
			wantErr:  assert.NoError,
		},
		{
			name: "error without data",
			fields: fields{
				makeClient: func(test *testing.T) Client {
					client, err := NewClient(opts.StorageAddress, "database", "collection")
					require.NoError(test, err)

					return client
				},
				keyField: CodeLinkField,
			},
			prepare: func(test *testing.T, getter LinkGetter) {
				_, err := getter.Client.
					Collection().
					DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)
			},
			args:     args{"code"},
			wantLink: entities.Link{},
			wantErr: func(test assert.TestingT, err error, args ...interface{}) bool {
				return assert.Equal(test, sql.ErrNoRows, err, args)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := data.fields.makeClient(test)
			getter := LinkGetter{
				Client:   client,
				KeyField: data.fields.keyField,
			}
			data.prepare(test, getter)

			gotLink, gotErr := getter.GetLink(data.args.query)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
