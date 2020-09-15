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
			[]mongo.IndexModel{
				mongo.IndexModel{
					Keys:    bson.D{{Key: "code", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
				mongo.IndexModel{
					Keys:    bson.D{{Key: "url", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
			},
			options.CreateIndexes(),
		)
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to create indexes in MongoDB")
	}

	return Client{innerClient: innerClient}, nil
}
