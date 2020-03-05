package generators

type chunkedCounter struct {
	step    uint64
	current uint64
	final   uint64
}

func newChunkedCounter(step uint64) chunkedCounter {
	return chunkedCounter{step: step, current: 0, final: 0}
}

func (counter chunkedCounter) isOver() bool {
	return counter.current == counter.final
}

func (counter *chunkedCounter) increase() (previous uint64) {
	previous = counter.current
	counter.current++

	return previous
}

func (counter *chunkedCounter) reset(initial uint64) {
	counter.current = initial
	counter.final = initial + counter.step
}
