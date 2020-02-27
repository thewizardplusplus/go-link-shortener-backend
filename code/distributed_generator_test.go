package code

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewDistributedGenerator(test *testing.T) {
	distributedCounters := []DistributedCounter{
		new(MockDistributedCounter),
		new(MockDistributedCounter),
	}
	randomSource := func(maximum int) int { panic("not implemented") }
	formatter := func(code uint64) string { panic("not implemented") }
	got :=
		NewDistributedGenerator(23, distributedCounters, randomSource, formatter)

	for _, distributedCounter := range distributedCounters {
		mock.AssertExpectationsForObjects(test, distributedCounter)
	}
	require.NotNil(test, got)
	assert.Equal(test, chunkedCounter{step: 23}, got.counter)
	assert.Equal(test, distributedCounters, got.distributedCounters)
	assert.Equal(test, getPointer(randomSource), getPointer(got.randomSource))
	assert.Equal(test, getPointer(formatter), getPointer(got.formatter))
}

func TestDistributedGenerator_GenerateCode(test *testing.T) {
	type fields struct {
		counter             chunkedCounter
		distributedCounters []DistributedCounter
		randomSource        RandomSource
		formatter           Formatter
	}

	for _, data := range []struct {
		name        string
		fields      fields
		wantCounter chunkedCounter
		wantCode    string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "success with incrementing",
			fields: fields{
				counter: chunkedCounter{step: 23, current: 42, final: 65},
				distributedCounters: []DistributedCounter{
					new(MockDistributedCounter),
					new(MockDistributedCounter),
				},
				randomSource: func(maximum int) int { panic("not implemented") },
				formatter:    func(code uint64) string { return fmt.Sprintf("[%d]", code) },
			},
			wantCounter: chunkedCounter{step: 23, current: 43, final: 65},
			wantCode:    "[42]",
			wantErr:     assert.NoError,
		},
		{
			name: "success with resetting",
			fields: fields{
				counter: chunkedCounter{step: 23, current: 65, final: 65},
				distributedCounters: func() []DistributedCounter {
					firstCounter := new(MockDistributedCounter)

					secondCounter := new(MockDistributedCounter)
					secondCounter.On("NextCountChunk").Return(uint64(100), nil)

					return []DistributedCounter{firstCounter, secondCounter}
				}(),
				randomSource: rand.New(rand.NewSource(1)).Intn,
				formatter:    func(code uint64) string { return fmt.Sprintf("[%d]", code) },
			},
			wantCounter: chunkedCounter{step: 23, current: 101, final: 123},
			wantCode:    "[100]",
			wantErr:     assert.NoError,
		},
		{
			name: "error with resetting",
			fields: fields{
				counter: chunkedCounter{step: 23, current: 65, final: 65},
				distributedCounters: func() []DistributedCounter {
					firstCounter := new(MockDistributedCounter)

					secondCounter := new(MockDistributedCounter)
					secondCounter.On("NextCountChunk").Return(uint64(0), iotest.ErrTimeout)

					return []DistributedCounter{firstCounter, secondCounter}
				}(),
				randomSource: rand.New(rand.NewSource(1)).Intn,
				formatter:    func(code uint64) string { panic("not implemented") },
			},
			wantCounter: chunkedCounter{step: 23, current: 65, final: 65},
			wantCode:    "",
			wantErr:     assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			generator := &DistributedGenerator{
				counter:             data.fields.counter,
				distributedCounters: data.fields.distributedCounters,
				randomSource:        data.fields.randomSource,
				formatter:           data.fields.formatter,
			}
			gotCode, gotErr := generator.GenerateCode()

			for _, distributedCounter := range data.fields.distributedCounters {
				mock.AssertExpectationsForObjects(test, distributedCounter)
			}
			assert.Equal(test, data.wantCounter, generator.counter)
			assert.Equal(test, data.wantCode, gotCode)
			data.wantErr(test, gotErr)
		})
	}
}

func TestDistributedGenerator_resetCounter(test *testing.T) {
	type fields struct {
		counter             chunkedCounter
		distributedCounters []DistributedCounter
		randomSource        RandomSource
		formatter           Formatter
	}

	for _, data := range []struct {
		name        string
		fields      fields
		wantCounter chunkedCounter
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				counter: chunkedCounter{step: 23, current: 42, final: 65},
				distributedCounters: func() []DistributedCounter {
					firstCounter := new(MockDistributedCounter)

					secondCounter := new(MockDistributedCounter)
					secondCounter.On("NextCountChunk").Return(uint64(100), nil)

					return []DistributedCounter{firstCounter, secondCounter}
				}(),
				randomSource: rand.New(rand.NewSource(1)).Intn,
				formatter:    func(code uint64) string { panic("not implemented") },
			},
			wantCounter: chunkedCounter{step: 23, current: 100, final: 123},
			wantErr:     assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				counter: chunkedCounter{step: 23, current: 42, final: 65},
				distributedCounters: func() []DistributedCounter {
					firstCounter := new(MockDistributedCounter)

					secondCounter := new(MockDistributedCounter)
					secondCounter.On("NextCountChunk").Return(uint64(0), iotest.ErrTimeout)

					return []DistributedCounter{firstCounter, secondCounter}
				}(),
				randomSource: rand.New(rand.NewSource(1)).Intn,
				formatter:    func(code uint64) string { panic("not implemented") },
			},
			wantCounter: chunkedCounter{step: 23, current: 42, final: 65},
			wantErr:     assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			generator := &DistributedGenerator{
				counter:             data.fields.counter,
				distributedCounters: data.fields.distributedCounters,
				randomSource:        data.fields.randomSource,
				formatter:           data.fields.formatter,
			}
			gotErr := generator.resetCounter()

			for _, distributedCounter := range data.fields.distributedCounters {
				mock.AssertExpectationsForObjects(test, distributedCounter)
			}
			assert.Equal(test, data.wantCounter, generator.counter)
			data.wantErr(test, gotErr)
		})
	}
}

func TestDistributedGenerator_selectCounter(test *testing.T) {
	type markedDistributedCounter struct {
		MockDistributedCounter

		ID int
	}

	distributedCounters := []DistributedCounter{
		&markedDistributedCounter{ID: 1},
		&markedDistributedCounter{ID: 2},
	}
	generator := &DistributedGenerator{
		counter:             chunkedCounter{step: 23},
		distributedCounters: distributedCounters,
		randomSource:        rand.New(rand.NewSource(1)).Intn,
		formatter:           func(code uint64) string { panic("not implemented") },
	}
	got := generator.selectCounter()

	for _, distributedCounter := range distributedCounters {
		mock.AssertExpectationsForObjects(test, distributedCounter)
	}
	assert.Equal(test, &markedDistributedCounter{ID: 2}, got)
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}
