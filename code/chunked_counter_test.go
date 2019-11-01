package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newChunkedCounter(test *testing.T) {
	got := newChunkedCounter(23)
	assert.Equal(test, chunkedCounter{step: 23}, got)
}

func Test_chunkedCounter_isOver(test *testing.T) {
	type fields struct {
		step    uint64
		current uint64
		final   uint64
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			counter := chunkedCounter{
				step:    data.fields.step,
				current: data.fields.current,
				final:   data.fields.final,
			}
			got := counter.isOver()

			assert.Equal(test, data.want, got)
		})
	}
}
