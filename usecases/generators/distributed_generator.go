package generators

import (
	"sync"

	"github.com/pkg/errors"
)

//go:generate mockery -name=DistributedCounter -inpkg -case=underscore -testonly

// DistributedCounter ...
type DistributedCounter interface {
	NextCountChunk() (uint64, error)
}

// RandomSource ...
//
// It should NOT be concurrently safe, because it'll call only under lock.
type RandomSource func(maximum int) int

// Formatter ...
type Formatter func(code uint64) string

// DistributedGenerator ...
type DistributedGenerator struct {
	locker              sync.Mutex
	counter             chunkedCounter
	distributedCounters []DistributedCounter
	randomSource        RandomSource
	formatter           Formatter
}

// NewDistributedGenerator ...
func NewDistributedGenerator(
	chunkSize uint64,
	distributedCounters []DistributedCounter,
	randomSource RandomSource,
	formatter Formatter,
) *DistributedGenerator {
	return &DistributedGenerator{
		counter:             newChunkedCounter(chunkSize),
		distributedCounters: distributedCounters,
		randomSource:        randomSource,
		formatter:           formatter,
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
	return generator.formatter(counter), nil
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
