package counters

//go:generate mockery -name=DistributedCounter -inpkg -case=underscore -testonly

// DistributedCounter ...
type DistributedCounter interface {
	NextCountChunk() (uint64, error)
}

// RandomSource ...
type RandomSource func(maximum int) int

// CounterGroup ...
type CounterGroup struct {
	DistributedCounters []DistributedCounter
	RandomSource        RandomSource
}

// SelectCounter ...
func (group CounterGroup) SelectCounter() DistributedCounter {
	index := group.RandomSource(len(group.DistributedCounters))
	return group.DistributedCounters[index]
}
