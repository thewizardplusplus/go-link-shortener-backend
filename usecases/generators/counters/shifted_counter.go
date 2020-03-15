package counters

// ShiftedCounter ...
type ShiftedCounter struct {
	DistributedCounter DistributedCounter
	Offset             uint64
}

// NextCountChunk ...
func (counter ShiftedCounter) NextCountChunk() (uint64, error) {
	countChunk, err := counter.DistributedCounter.NextCountChunk()
	if err != nil {
		return 0, err
	}

	return countChunk + counter.Offset, nil
}
