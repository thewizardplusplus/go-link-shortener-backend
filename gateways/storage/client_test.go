// +build integration

package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Index struct {
	Name      string `bson:"name"`
	Namespace string `bson:"ns"`
	Key       bson.M `bson:"key"`
	Unique    bool   `bson:"unique"`
}

func TestNewClient(test *testing.T) {
	type args struct {
		uri        string
		database   string
		collection string
	}

	for _, data := range []struct {
		name        string
		args        args
		wantClient  require.ValueAssertionFunc
		wantIndexes []Index
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				uri:        "mongodb://localhost:27017",
				database:   "database",
				collection: "collection",
			},
			wantClient: require.NotNil,
			wantIndexes: []Index{
				{
					Name:      "_id_",
					Namespace: "database.collection",
					Key:       bson.M{"_id": int32(1)},
					Unique:    false,
				},
				{
					Name:      CodeLinkField + "_1",
					Namespace: "database.collection",
					Key:       bson.M{CodeLinkField: int32(1)},
					Unique:    true,
				},
				{
					Name:      URLLinkField + "_1",
					Namespace: "database.collection",
					Key:       bson.M{URLLinkField: int32(1)},
					Unique:    true,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				uri:        ":",
				database:   "database",
				collection: "collection",
			},
			wantClient:  require.Nil,
			wantIndexes: nil,
			wantErr:     assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotClient, gotErr :=
				NewClient(data.args.uri, data.args.database, data.args.collection)

			data.wantClient(test, gotClient.innerClient)
			data.wantErr(test, gotErr)
			if gotErr != nil {
				return
			}

			cursor, err := gotClient.
				Collection().
				Indexes().
				List(context.Background(), options.ListIndexes())
			require.NoError(test, err)

			var indexes []Index
			err = cursor.All(context.Background(), &indexes)
			require.NoError(test, err)

			assert.ElementsMatch(test, data.wantIndexes, indexes)
			assert.Equal(test, data.args.database, gotClient.database)
			assert.Equal(test, data.args.collection, gotClient.collection)
		})
	}
}
