package code

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewDistributedGenerator(test *testing.T) {
	distributedCounters := []DistributedCounter{
		new(MockDistributedCounter),
		new(MockDistributedCounter),
	}
	got := NewDistributedGenerator(23, distributedCounters, rand.Intn)

	for _, distributedCounter := range distributedCounters {
		mock.AssertExpectationsForObjects(test, distributedCounter)
	}
	require.NotNil(test, got)
	assert.Equal(test, chunkedCounter{step: 23}, got.counter)
	assert.Equal(test, distributedCounters, got.distributedCounters)
	assert.Equal(test, getPointer(rand.Intn), getPointer(got.randomSource))
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}
