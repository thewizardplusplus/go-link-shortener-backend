package counters

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
