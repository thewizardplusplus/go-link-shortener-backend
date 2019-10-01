package storage

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Client ...
type Client struct {
	innerClient *mongo.Client
}
