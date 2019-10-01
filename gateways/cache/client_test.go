package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(test *testing.T) {
	client := NewClient("localhost:6379")

	require.NotNil(test, client.innerClient)
	assert.Equal(test, "localhost:6379", client.innerClient.Options().Addr)
}
