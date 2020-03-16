package transformers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLinear(test *testing.T) {
	type fields struct {
		options []LinearOption
	}
	type args struct {
		countChunk uint64
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   uint64
	}{
		{
			name: "without options",
			fields: fields{
				options: nil,
			},
			args: args{5},
			want: 5,
		},
		{
			name: "with the factor",
			fields: fields{
				options: []LinearOption{WithFactor(2)},
			},
			args: args{5},
			want: 10,
		},
		{
			name: "with the offset",
			fields: fields{
				options: []LinearOption{WithOffset(3)},
			},
			args: args{5},
			want: 8,
		},
		{
			name: "with the factor and the offset",
			fields: fields{
				options: []LinearOption{WithFactor(2), WithOffset(3)},
			},
			args: args{5},
			want: 13,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			transformer := NewLinear(data.fields.options...)
			got := transformer(data.args.countChunk)

			assert.Equal(test, data.want, got)
		})
	}
}
