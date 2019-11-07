package counter

import (
	"context"

	"github.com/pkg/errors"
)

// Counter ...
type Counter struct {
	Client Client
	Name   string
}

// NextCountChunk ...
func (counter Counter) NextCountChunk() (uint64, error) {
	context := context.Background()
	response, err := counter.Client.innerClient.Put(context, counter.Name, "")
	if err != nil {
		return 0, errors.Wrap(err, "unable to update the counter")
	}

	return uint64(response.Header.Revision), nil
}
