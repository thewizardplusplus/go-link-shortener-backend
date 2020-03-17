package generators

// nolint: lll
import (
	"sync"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators/counters"
)

// nolint: lll
//go:generate mockery -name=DistributedCounterGroup -inpkg -case=underscore -testonly

// DistributedCounterGroup ...
type DistributedCounterGroup interface {
	SelectCounter() counters.DistributedCounter
}

// Formatter ...
type Formatter func(code uint64) string

// DistributedGenerator ...
type DistributedGenerator struct {
	locker              sync.Mutex
	counter             counters.ChunkedCounter
	distributedCounters DistributedCounterGroup
	formatter           Formatter
}

// NewDistributedGenerator ...
func NewDistributedGenerator(
	chunkSize uint64,
	distributedCounters DistributedCounterGroup,
	formatter Formatter,
) *DistributedGenerator {
	return &DistributedGenerator{
		counter:             counters.NewChunkedCounter(chunkSize),
		distributedCounters: distributedCounters,
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
	countChunk, err :=
		generator.distributedCounters.SelectCounter().NextCountChunk()
	if err != nil {
		return errors.Wrap(err, "unable to get the next count chunk")
	}

	generator.counter.Reset(countChunk)
	return nil
}
