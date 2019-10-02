package storage

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener/entities"
	"go.mongodb.org/mongo-driver/bson"
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
	if err != nil {
		return entities.Link{},
			errors.Wrap(err, "unable to get the link from MongoDB")
	}

	return link, nil
}
