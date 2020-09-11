// +build integration

package storage

import (
	"context"
	"testing"

	"github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLinkSetter_SetLink(test *testing.T) {
	// nolint: lll
	type options struct {
		StorageAddress string `env:"STORAGE_ADDRESS" envDefault:"mongodb://localhost:27017"`
	}
	type fields struct {
		makeClient func(test *testing.T) Client
		database   string
		collection string
	}
	type args struct {
		link entities.Link
	}

	var opts options
	err := env.Parse(&opts)
	require.NoError(test, err)

	for _, data := range []struct {
		name    string
		fields  fields
		prepare func(test *testing.T, setter LinkSetter)
		args    args
		wantErr assert.ErrorAssertionFunc
		check   func(test *testing.T, setter LinkSetter)
	}{
		{
			name: "success with creating",
			fields: fields{
				makeClient: func(test *testing.T) Client {
					client, err := NewClient(opts.StorageAddress)
					require.NoError(test, err)

					return client
				},
				database:   "database",
				collection: "collection",
			},
			prepare: func(test *testing.T, setter LinkSetter) {
				_, err := setter.Client.innerClient.
					Database(setter.Database).
					Collection(setter.Collection).
					DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)
			},
			args: args{
				link: entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, setter LinkSetter) {
				cursor, err := setter.Client.innerClient.
					Database(setter.Database).
					Collection(setter.Collection).
					Find(context.Background(), bson.M{"url": "url"})
				require.NoError(test, err)

				var links []entities.Link
				err = cursor.All(context.Background(), &links)
				require.NoError(test, err)

				assert.Equal(test, []entities.Link{{Code: "code", URL: "url"}}, links)
			},
		},
		{
			name: "success with updating",
			fields: fields{
				makeClient: func(test *testing.T) Client {
					client, err := NewClient(opts.StorageAddress)
					require.NoError(test, err)

					return client
				},
				database:   "database",
				collection: "collection",
			},
			prepare: func(test *testing.T, setter LinkSetter) {
				_, err := setter.Client.innerClient.
					Database(setter.Database).
					Collection(setter.Collection).
					DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)

				_, err = setter.Client.innerClient.
					Database(setter.Database).
					Collection(setter.Collection).
					InsertOne(context.Background(), entities.Link{Code: "code #1", URL: "url"})
				require.NoError(test, err)
			},
			args: args{
				link: entities.Link{Code: "code #2", URL: "url"},
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, setter LinkSetter) {
				cursor, err := setter.Client.innerClient.
					Database(setter.Database).
					Collection(setter.Collection).
					Find(context.Background(), bson.M{"url": "url"})
				require.NoError(test, err)

				var links []entities.Link
				err = cursor.All(context.Background(), &links)
				require.NoError(test, err)

				assert.Equal(test, []entities.Link{{Code: "code #1", URL: "url"}}, links)
			},
		},
		{
			name: "success with skipping",
			fields: fields{
				makeClient: func(test *testing.T) Client {
					client, err := NewClient(opts.StorageAddress)
					require.NoError(test, err)

					return client
				},
				database:   "database",
				collection: "collection",
			},
			prepare: func(test *testing.T, setter LinkSetter) {
				_, err := setter.Client.innerClient.
					Database(setter.Database).
					Collection(setter.Collection).
					DeleteMany(context.Background(), bson.M{})
				require.NoError(test, err)

				_, err = setter.Client.innerClient.
					Database(setter.Database).
					Collection(setter.Collection).
					InsertOne(context.Background(), entities.Link{Code: "code", URL: "url"})
				require.NoError(test, err)
			},
			args: args{
				link: entities.Link{Code: "code", URL: "url"},
			},
			wantErr: assert.NoError,
			check: func(test *testing.T, setter LinkSetter) {
				cursor, err := setter.Client.innerClient.
					Database(setter.Database).
					Collection(setter.Collection).
					Find(context.Background(), bson.M{"url": "url"})
				require.NoError(test, err)

				var links []entities.Link
				err = cursor.All(context.Background(), &links)
				require.NoError(test, err)

				assert.Equal(test, []entities.Link{{Code: "code", URL: "url"}}, links)
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			client := data.fields.makeClient(test)
			setter := LinkSetter{
				Client:     client,
				Database:   data.fields.database,
				Collection: data.fields.collection,
			}
			data.prepare(test, setter)

			gotErr := setter.SetLink(data.args.link)

			data.wantErr(test, gotErr)
			data.check(test, setter)
		})
	}
}
