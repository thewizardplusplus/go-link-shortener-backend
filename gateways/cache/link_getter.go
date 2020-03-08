package cache

import (
	"database/sql"
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
)

// LinkGetter ...
type LinkGetter struct {
	Client Client
}

// GetLink ...
func (getter LinkGetter) GetLink(query string) (entities.Link, error) {
	data, err := getter.Client.innerClient.Get(query).Result()
	switch err {
	case nil:
	case redis.Nil:
		return entities.Link{}, sql.ErrNoRows
	default:
		return entities.Link{}, errors.Wrap(err, "unable to get the link from Redis")
	}

	var link entities.Link
	if err := json.Unmarshal([]byte(data), &link); err != nil {
		return entities.Link{},
			errors.Wrap(err, "unable to unmarshal the link from Redis")
	}

	return link, nil
}
