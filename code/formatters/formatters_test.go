package formatters

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInBase10(test *testing.T) {
	type args struct {
		code uint64
	}

	for _, data := range []struct {
		name string
		args args
		want string
	}{
		{
			name: "regular number",
			args: args{123456789},
			want: "123456789",
		},
		{
			name: "minimal number",
			args: args{0},
			want: "0",
		},
		{
			name: "maximal number",
			args: args{math.MaxUint64},
			want: "18446744073709551615",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := InBase10(data.args.code)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestInBase62(test *testing.T) {
	type args struct {
		code uint64
	}

	for _, data := range []struct {
		name string
		args args
		want string
	}{
		{
			name: "regular number",
			args: args{123456789},
			want: "8m0Kx",
		},
		{
			name: "minimal number",
			args: args{0},
			want: "0",
		},
		{
			name: "maximal number",
			args: args{math.MaxUint64},
			want: "lYGhA16ahyf",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := InBase62(data.args.code)

			assert.Equal(test, data.want, got)
		})
	}
}
