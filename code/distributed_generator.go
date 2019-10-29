package code

import (
	"sync"

	"github.com/pkg/errors"
)

// RandomSource ...
//
// It should NOT be concurrently safe, because it'll call only under lock.
type RandomSource func(maximum int) int

// DistributedCounter ...
type DistributedCounter interface {
	NextCountChunk(counterName string) (uint64, error)
}

// DistributedGenerator ...
type DistributedGenerator struct {
	CountersNames      []string
	RandomSource       RandomSource
	DistributedCounter DistributedCounter
	ChunkSize          uint64

	locker  sync.Mutex
	counter uint64
	limit   uint64
}

func (generator *DistributedGenerator) resetCounter() error {
	countChunk, err := generator.DistributedCounter.
		NextCountChunk(generator.counterName())
	if err != nil {
		return errors.Wrap(err, "unable to get the next count chunk")
	}

	generator.counter = countChunk
	generator.limit = countChunk + generator.ChunkSize

	return nil
}

func (generator *DistributedGenerator) counterName() string {
	index := generator.RandomSource(len(generator.CountersNames))
	return generator.CountersNames[index]
}
