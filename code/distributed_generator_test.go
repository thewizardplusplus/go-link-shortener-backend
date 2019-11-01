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

func TestDistributedGenerator_selectCounter(test *testing.T) {
	type markedDistributedCounter struct {
		MockDistributedCounter

		ID int
	}

	distributedCounters := []DistributedCounter{
		&markedDistributedCounter{ID: 1},
		&markedDistributedCounter{ID: 2},
	}
	randomSource := rand.New(rand.NewSource(1))
	generator := &DistributedGenerator{
		counter:             chunkedCounter{step: 23},
		distributedCounters: distributedCounters,
		randomSource:        randomSource.Intn,
	}
	got := generator.selectCounter()

	for _, distributedCounter := range distributedCounters {
		mock.AssertExpectationsForObjects(test, distributedCounter)
	}
	assert.Equal(test, &markedDistributedCounter{ID: 2}, got)
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}
