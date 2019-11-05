package counter

import (
	"go.etcd.io/etcd/clientv3"
)

// Client ...
type Client struct {
	innerClient *clientv3.Client
}
