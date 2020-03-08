package storage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// LinkGetter ...
type LinkGetter struct {
	Client     Client
	Database   string
	Collection string
	KeyField   string
}

// GetLink ...
func (getter LinkGetter) GetLink(query string) (entities.Link, error) {
	var link entities.Link
	err := getter.Client.innerClient.
		Database(getter.Database).
		Collection(getter.Collection).
		FindOne(context.Background(), bson.M{getter.KeyField: query}).
		Decode(&link)
	switch err {
	case nil:
		return link, nil
	case mongo.ErrNoDocuments:
		return entities.Link{}, sql.ErrNoRows
	default:
		return entities.Link{},
			errors.Wrap(err, "unable to get the link from MongoDB")
	}
}
