package counters

// MultipliedCounter ...
type MultipliedCounter struct {
	DistributedCounter DistributedCounter
	Factor             uint64
}

// NextCountChunk ...
func (counter MultipliedCounter) NextCountChunk() (uint64, error) {
	countChunk, err := counter.DistributedCounter.NextCountChunk()
	if err != nil {
		return 0, err
	}

	return countChunk * counter.Factor, nil
}
