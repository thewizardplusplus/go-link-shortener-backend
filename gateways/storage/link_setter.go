package storage

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// LinkSetter ...
type LinkSetter struct {
	Client     Client
	Database   string
	Collection string
}

// SetLink ...
func (setter LinkSetter) SetLink(link entities.Link) error {
	_, err := setter.Client.innerClient.
		Database(setter.Database).
		Collection(setter.Collection).
		InsertOne(context.Background(), link)
	if err != nil {
		return errors.Wrap(err, "unable to set the link in MongoDB")
	}

	return nil
}
