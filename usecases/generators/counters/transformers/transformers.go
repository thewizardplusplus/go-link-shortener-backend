package transformers

// LinearConfig ...
type LinearConfig struct {
	factor uint64
	offset uint64
}

// LinearOption ...
type LinearOption func(config *LinearConfig)

// WithFactor ...
func WithFactor(factor uint64) LinearOption {
	return func(config *LinearConfig) { config.factor = factor }
}

// WithOffset ...
func WithOffset(offset uint64) LinearOption {
	return func(config *LinearConfig) { config.offset = offset }
}

// Linear ...
func Linear(countChunk uint64, options ...LinearOption) uint64 {
	config := LinearConfig{
		factor: 1,
		offset: 0,
	}
	for _, option := range options {
		option(&config)
	}

	return countChunk*config.factor + config.offset
}
