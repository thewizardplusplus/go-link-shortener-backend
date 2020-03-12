package counters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChunkedCounter(test *testing.T) {
	got := NewChunkedCounter(23)

	assert.Equal(test, ChunkedCounter{step: 23}, got)
}

func TestChunkedCounter_IsOver(test *testing.T) {
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
		{
			name: "current less than final",
			fields: fields{
				step:    23,
				current: 42,
				final:   65,
			},
			want: false,
		},
		{
			name: "current equal to final",
			fields: fields{
				step:    23,
				current: 65,
				final:   65,
			},
			want: true,
		},
		{
			name: "current greater than final",
			fields: fields{
				step:    23,
				current: 100,
				final:   65,
			},
			want: true,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			counter := ChunkedCounter{
				step:    data.fields.step,
				current: data.fields.current,
				final:   data.fields.final,
			}
			got := counter.IsOver()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestChunkedCounter_Increase(test *testing.T) {
	counter := ChunkedCounter{current: 23}
	previous := counter.Increase()

	assert.Equal(test, ChunkedCounter{current: 24}, counter)
	assert.Equal(test, uint64(23), previous)
}

func TestChunkedCounter_Reset(test *testing.T) {
	counter := ChunkedCounter{step: 23, current: 42, final: 65}
	counter.Reset(100)

	assert.Equal(test, ChunkedCounter{step: 23, current: 100, final: 123}, counter)
}
