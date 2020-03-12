package generators

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/counters"
)

// RandomSource ...
//
// It should NOT be concurrently safe, because it'll call only under lock.
type RandomSource func(maximum int) int

// Formatter ...
type Formatter func(code uint64) string

// DistributedGenerator ...
type DistributedGenerator struct {
	locker              sync.Mutex
	counter             counters.ChunkedCounter
	distributedCounters []counters.DistributedCounter
	randomSource        RandomSource
	formatter           Formatter
}

// NewDistributedGenerator ...
func NewDistributedGenerator(
	chunkSize uint64,
	distributedCounters []counters.DistributedCounter,
	randomSource RandomSource,
	formatter Formatter,
) *DistributedGenerator {
	return &DistributedGenerator{
		counter:             counters.NewChunkedCounter(chunkSize),
		distributedCounters: distributedCounters,
		randomSource:        randomSource,
		formatter:           formatter,
	}
}

// GenerateCode ...
func (generator *DistributedGenerator) GenerateCode() (string, error) {
	generator.locker.Lock()
	defer generator.locker.Unlock()

	if generator.counter.IsOver() {
		if err := generator.resetCounter(); err != nil {
			return "", errors.Wrap(err, "unable to reset the counter")
		}
	}

	counter := generator.counter.Increase()
	return generator.formatter(counter), nil
}

func (generator *DistributedGenerator) resetCounter() error {
	countChunk, err := generator.selectCounter().NextCountChunk()
	if err != nil {
		return errors.Wrap(err, "unable to get the next count chunk")
	}

	generator.counter.Reset(countChunk)
	return nil
}

func (generator *DistributedGenerator) selectCounter() counters.DistributedCounter {
	index := generator.randomSource(len(generator.distributedCounters))
	return generator.distributedCounters[index]
}
