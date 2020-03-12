package counters

// DistributedCounter ...
type DistributedCounter interface {
	NextCountChunk() (uint64, error)
}
