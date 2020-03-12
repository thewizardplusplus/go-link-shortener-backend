package counters

// ChunkedCounter ...
type ChunkedCounter struct {
	step    uint64
	current uint64
	final   uint64
}

// NewChunkedCounter ...
func NewChunkedCounter(step uint64) ChunkedCounter {
	return ChunkedCounter{step: step, current: 0, final: 0}
}

// IsOver ...
func (counter ChunkedCounter) IsOver() bool {
	return counter.current >= counter.final
}

// Increase ...
func (counter *ChunkedCounter) Increase() (previous uint64) {
	previous = counter.current
	counter.current++

	return previous
}

// Reset ...
func (counter *ChunkedCounter) Reset(initial uint64) {
	counter.current = initial
	counter.final = initial + counter.step
}
