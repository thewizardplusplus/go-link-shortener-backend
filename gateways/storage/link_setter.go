package storage

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LinkSetter ...
type LinkSetter struct {
	Client Client
}

// SetLink ...
func (setter LinkSetter) SetLink(link entities.Link) error {
	// by the time of setting the database may already have a link created
	// in another thread; therefore, to avoid duplicates, we don't insert
	// but update in the upsert mode; a link code is always unique, so we search
	// by a link URL
	_, err := setter.Client.innerClient.
		Database(setter.Client.database).
		Collection(setter.Client.collection).
		UpdateOne(
			context.Background(),
			bson.M{"url": link.URL},
			bson.M{"$setOnInsert": bson.M{"code": link.Code}},
			options.Update().SetUpsert(true),
		)
	if err != nil {
		return errors.Wrap(err, "unable to set the link in MongoDB")
	}

	return nil
}
