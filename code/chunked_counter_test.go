package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newChunkedCounter(test *testing.T) {
	got := newChunkedCounter(23)
	assert.Equal(test, chunkedCounter{step: 23}, got)
}
