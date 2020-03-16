package counters

//go:generate mockery -name=DistributedCounter -inpkg -case=underscore -testonly

// DistributedCounter ...
type DistributedCounter interface {
	NextCountChunk() (uint64, error)
}

// Transformer ...
type Transformer func(countChunk uint64) uint64

// TransformedCounter ...
type TransformedCounter struct {
	DistributedCounter DistributedCounter
	Transformer        Transformer
}

// NextCountChunk ...
func (counter TransformedCounter) NextCountChunk() (uint64, error) {
	countChunk, err := counter.DistributedCounter.NextCountChunk()
	if err != nil {
		return 0, err
	}

	return counter.Transformer(countChunk), nil
}
