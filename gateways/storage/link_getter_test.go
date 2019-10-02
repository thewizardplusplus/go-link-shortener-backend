// +build integration

package storage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener/entities"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	address string
)

// nolint: gochecknoinits
func init() {
	var ok bool
	if address, ok = os.LookupEnv("MONGODB_URL"); !ok {
		address = "mongodb://localhost:27017"
	}
}

func TestLinkGetter_GetLink(test *testing.T) {
	type fields struct {
		makeClient func(test *testing.T) Client
		database   string
		collection string
		keyField   string
	}
	type args struct {
		query string
	}

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
					client, err := NewClient(address)
					require.NoError(test, err)

					return client
				},
				database:   "database",
				collection: "collection",
				keyField:   "code",
			},
			prepare: func(test *testing.T, getter LinkGetter) {
				_, err := getter.Client.innerClient.
					Database(getter.Database).
					Collection(getter.Collection).
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
					client, err := NewClient(address)
					require.NoError(test, err)

					return client
				},
				database:   "database",
				collection: "collection",
				keyField:   "code",
			},
			prepare: func(test *testing.T, getter LinkGetter) {
				_, err := getter.Client.innerClient.
					Database(getter.Database).
					Collection(getter.Collection).
					DeleteMany(context.Background(), bson.M{"code": "code"})
				require.NoError(test, err)
			},
			args:     args{"code"},
			wantLink: entities.Link{},
			wantErr:  assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := data.fields.makeClient(test)
			getter := LinkGetter{
				Client:     client,
				Database:   data.fields.database,
				Collection: data.fields.collection,
				KeyField:   data.fields.keyField,
			}
			data.prepare(test, getter)

			gotLink, gotErr := getter.GetLink(data.args.query)

			assert.Equal(test, data.wantLink, gotLink)
			data.wantErr(test, gotErr)
		})
	}
}
