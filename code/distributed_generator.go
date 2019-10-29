package code

import (
	"sync"
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

func (generator *DistributedGenerator) counterName() string {
	index := generator.RandomSource(len(generator.CountersNames))
	return generator.CountersNames[index]
}
