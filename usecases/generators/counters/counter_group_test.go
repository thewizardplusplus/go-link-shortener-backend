package counters

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MarkedDistributedCounter struct {
	MockDistributedCounter

	ID int
}

func NewMarkedDistributedCounter(id int) *MarkedDistributedCounter {
	return &MarkedDistributedCounter{ID: id}
}

func TestCounterGroup_SelectCounter(test *testing.T) {
	distributedCounters := []DistributedCounter{
		NewMarkedDistributedCounter(1),
		NewMarkedDistributedCounter(2),
	}
	group := &CounterGroup{
		DistributedCounters: distributedCounters,
		RandomSource:        rand.New(rand.NewSource(1)).Intn,
	}
	got := group.SelectCounter()

	for _, distributedCounter := range distributedCounters {
		mock.AssertExpectationsForObjects(test, distributedCounter)
	}
	assert.Equal(test, NewMarkedDistributedCounter(2), got)
}
