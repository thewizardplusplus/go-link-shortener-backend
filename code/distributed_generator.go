package code

import (
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

// RandomSource ...
//
// It should NOT be concurrently safe, because it'll call only under lock.
type RandomSource func(maximum int) int

// DistributedCounter ...
type DistributedCounter interface {
	NextCountChunk() (uint64, error)
}

// DistributedGenerator ...
type DistributedGenerator struct {
	RandomSource        RandomSource
	DistributedCounters []DistributedCounter
	ChunkSize           uint64

	locker  sync.Mutex
	counter uint64
	limit   uint64
}

// GenerateCode ...
func (generator *DistributedGenerator) GenerateCode() (string, error) {
	generator.locker.Lock()
	defer generator.locker.Unlock()

	if generator.counter == generator.limit {
		if err := generator.resetCounter(); err != nil {
			return "", errors.Wrap(err, "unable to reset the counter")
		}
	}

	counter := generator.counter
	generator.counter++

	return strconv.FormatUint(counter, 10), nil
}

func (generator *DistributedGenerator) resetCounter() error {
	countChunk, err := generator.selectCounter().NextCountChunk()
	if err != nil {
		return errors.Wrap(err, "unable to get the next count chunk")
	}

	generator.counter = countChunk
	generator.limit = countChunk + generator.ChunkSize

	return nil
}

func (generator *DistributedGenerator) selectCounter() DistributedCounter {
	index := generator.RandomSource(len(generator.DistributedCounters))
	return generator.DistributedCounters[index]
}
