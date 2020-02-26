package counter

import (
	"context"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

// Counter ...
type Counter struct {
	Client Client
	Name   string
}

// NextCountChunk ...
func (counter Counter) NextCountChunk() (uint64, error) {
	response, err := counter.Client.innerClient.
		Put(context.Background(), counter.Name, "", clientv3.WithPrevKV())
	if err != nil {
		return 0, errors.Wrap(err, "unable to update the counter")
	}
	if response.PrevKv == nil {
		return 0, errors.Wrap(err, "unable to get the previous counter")
	}

	return uint64(response.PrevKv.Version) + 1, nil
}
