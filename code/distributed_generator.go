package code

import (
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

// DistributedCounter ...
type DistributedCounter interface {
	NextCountChunk() (uint64, error)
}

// RandomSource ...
//
// It should NOT be concurrently safe, because it'll call only under lock.
type RandomSource func(maximum int) int

// DistributedGenerator ...
type DistributedGenerator struct {
	locker              sync.Mutex
	counter             chunkedCounter
	distributedCounters []DistributedCounter
	randomSource        RandomSource
}

// NewDistributedGenerator ...
func NewDistributedGenerator(
	chunkSize uint64,
	distributedCounters []DistributedCounter,
	randomSource RandomSource,
) *DistributedGenerator {
	return &DistributedGenerator{
		counter:             newChunkedCounter(chunkSize),
		distributedCounters: distributedCounters,
		randomSource:        randomSource,
	}
}

// GenerateCode ...
func (generator *DistributedGenerator) GenerateCode() (string, error) {
	generator.locker.Lock()
	defer generator.locker.Unlock()

	if generator.counter.isOver() {
		if err := generator.resetCounter(); err != nil {
			return "", errors.Wrap(err, "unable to reset the counter")
		}
	}

	counter := generator.counter.increase()
	return strconv.FormatUint(counter, 10), nil
}

func (generator *DistributedGenerator) resetCounter() error {
	countChunk, err := generator.selectCounter().NextCountChunk()
	if err != nil {
		return errors.Wrap(err, "unable to get the next count chunk")
	}

	generator.counter.reset(countChunk)
	return nil
}

func (generator *DistributedGenerator) selectCounter() DistributedCounter {
	index := generator.randomSource(len(generator.distributedCounters))
	return generator.distributedCounters[index]
}
