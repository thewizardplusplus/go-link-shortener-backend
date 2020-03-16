package transformers

// nolint: lll
import (
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators/counters"
)

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

// NewLinear ...
func NewLinear(options ...LinearOption) counters.Transformer {
	config := LinearConfig{
		factor: 1,
		offset: 0,
	}
	for _, option := range options {
		option(&config)
	}

	return func(countChunk uint64) uint64 {
		return countChunk*config.factor + config.offset
	}
}
