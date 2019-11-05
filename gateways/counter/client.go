package counter

import (
	"github.com/pkg/errors"
	"go.etcd.io/etcd/clientv3"
)

// Client ...
type Client struct {
	innerClient *clientv3.Client
}

// NewClient ...
func NewClient(url string) (Client, error) {
	innerClient, err := clientv3.New(clientv3.Config{Endpoints: []string{url}})
	if err != nil {
		return Client{}, errors.Wrap(err, "unable to connect to etcd")
	}

	return Client{innerClient: innerClient}, nil
}
