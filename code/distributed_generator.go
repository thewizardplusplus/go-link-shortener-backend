package code

// RandomSource ...
//
// It should NOT be concurrently safe, because it'll call only under lock.
type RandomSource func(maximum int) int

// DistributedCounter ...
type DistributedCounter interface {
	NextCountChunk(counterName string) (uint64, error)
}
