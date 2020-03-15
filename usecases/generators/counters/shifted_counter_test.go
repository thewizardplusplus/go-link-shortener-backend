package counters

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShiftedCounter_NextCountChunk(test *testing.T) {
	type fields struct {
		DistributedCounter DistributedCounter
		Offset             uint64
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
				Offset: 1000,
			},
			wantChunk: 1005,
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
				Offset: 1000,
			},
			wantChunk: 0,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			counter := ShiftedCounter{
				DistributedCounter: data.fields.DistributedCounter,
				Offset:             data.fields.Offset,
			}
			gotChunk, gotErr := counter.NextCountChunk()

			mock.AssertExpectationsForObjects(test, data.fields.DistributedCounter)
			assert.Equal(test, data.wantChunk, gotChunk)
			data.wantErr(test, gotErr)
		})
	}
}
