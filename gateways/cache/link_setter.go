package cache

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener/entities"
)

// KeyExtractor ...
type KeyExtractor func(link entities.Link) string

// LinkSetter ...
type LinkSetter struct {
	KeyExtractor KeyExtractor
	Client       Client
	Expiration   time.Duration
}

// SetLink ...
func (setter LinkSetter) SetLink(link entities.Link) error {
	data, err := json.Marshal(link)
	if err != nil {
		return errors.Wrap(err, "unable to marshal the link for Redis")
	}

	key := setter.KeyExtractor(link)
	if err := setter.Client.innerClient.
		Set(key, string(data), setter.Expiration).
		Err(); err != nil {
		return errors.Wrap(err, "unable to set the link in Redis")
	}

	return nil
}
