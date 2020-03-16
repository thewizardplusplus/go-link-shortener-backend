package counters

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransformedCounter_NextCountChunk(test *testing.T) {
	type fields struct {
		DistributedCounter DistributedCounter
		Transformer        Transformer
	}

	for _, data := range []struct {
		name      string
		fields    fields
		wantChunk uint64
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				DistributedCounter: func() DistributedCounter {
					counter := new(MockDistributedCounter)
					counter.On("NextCountChunk").Return(uint64(5), nil)

					return counter
				}(),
				Transformer: func(countChunk uint64) uint64 { return countChunk * 2 },
			},
			wantChunk: 10,
			wantErr:   assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				DistributedCounter: func() DistributedCounter {
					counter := new(MockDistributedCounter)
					counter.On("NextCountChunk").Return(uint64(0), iotest.ErrTimeout)

					return counter
				}(),
				Transformer: func(countChunk uint64) uint64 { panic("not implemented") },
			},
			wantChunk: 0,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			counter := TransformedCounter{
				DistributedCounter: data.fields.DistributedCounter,
				Transformer:        data.fields.Transformer,
			}
			gotChunk, gotErr := counter.NextCountChunk()

			mock.AssertExpectationsForObjects(test, data.fields.DistributedCounter)
			assert.Equal(test, data.wantChunk, gotChunk)
			data.wantErr(test, gotErr)
		})
	}
}
