package counters

//go:generate mockery -name=DistributedCounter -inpkg -case=underscore -testonly

// DistributedCounter ...
type DistributedCounter interface {
	NextCountChunk() (uint64, error)
}
