package code

type chunkedCounter struct {
	step    uint64
	current uint64
	final   uint64
}

func newChunkedCounter(step uint64) chunkedCounter {
	return chunkedCounter{step: step, current: 0, final: 0}
}
