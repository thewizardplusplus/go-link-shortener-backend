package storage

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client ...
type Client struct {
	innerClient *mongo.Client
	database    string
	collection  string
}

// NewClient ...
func NewClient(uri string, database string, collection string) (Client, error) {
	innerClient, err :=
		mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to connect to MongoDB")
	}

	_, err = innerClient.
		Database(database).
		Collection(collection).
		Indexes().
		CreateMany(
			context.Background(),
			[]mongo.IndexModel{makeUniqueIndex("code"), makeUniqueIndex("url")},
			options.CreateIndexes(),
		)
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to create indexes in MongoDB")
	}

	client := Client{
		innerClient: innerClient,
		database:    database,
		collection:  collection,
	}
	return client, nil
}

func makeUniqueIndex(key string) mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    bson.D{{Key: key, Value: 1}},
		Options: options.Index().SetUnique(true),
	}
}
