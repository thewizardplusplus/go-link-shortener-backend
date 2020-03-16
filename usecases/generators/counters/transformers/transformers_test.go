package transformers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinear(test *testing.T) {
	type args struct {
		countChunk uint64
		options    []LinearOption
	}

	for _, data := range []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "without options",
			args: args{
				countChunk: 5,
				options:    nil,
			},
			want: 5,
		},
		{
			name: "with the factor",
			args: args{
				countChunk: 5,
				options:    []LinearOption{WithFactor(2)},
			},
			want: 10,
		},
		{
			name: "with the offset",
			args: args{
				countChunk: 5,
				options:    []LinearOption{WithOffset(3)},
			},
			want: 8,
		},
		{
			name: "with the factor and the offset",
			args: args{
				countChunk: 5,
				options:    []LinearOption{WithFactor(2), WithOffset(3)},
			},
			want: 13,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := Linear(data.args.countChunk, data.args.options...)

			assert.Equal(test, data.want, got)
		})
	}
}
