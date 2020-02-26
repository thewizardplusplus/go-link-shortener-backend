package storage

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client ...
type Client struct {
	innerClient *mongo.Client
}

// NewClient ...
func NewClient(uri string) (Client, error) {
	options := options.Client().ApplyURI(uri)
	innerClient, err := mongo.Connect(context.Background(), options)
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to connect to MongoDB")
	}

	return Client{innerClient: innerClient}, nil
}
