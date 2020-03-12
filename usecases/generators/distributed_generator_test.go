package generators

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/counters"
)

type MarkedDistributedCounter struct {
	MockDistributedCounter

	ID int
}

func NewMarkedDistributedCounter(id int) *MarkedDistributedCounter {
	return &MarkedDistributedCounter{ID: id}
}

func TestNewDistributedGenerator(test *testing.T) {
	distributedCounters := []counters.DistributedCounter{
		NewMarkedDistributedCounter(1),
		NewMarkedDistributedCounter(2),
	}
	randomSource := func(maximum int) int { panic("not implemented") }
	formatter := func(code uint64) string { panic("not implemented") }
	got :=
		NewDistributedGenerator(23, distributedCounters, randomSource, formatter)

	for _, distributedCounter := range distributedCounters {
		mock.AssertExpectationsForObjects(test, distributedCounter)
	}
	require.NotNil(test, got)
	assert.Equal(test, counters.NewChunkedCounter(23), got.counter)
	assert.Equal(test, distributedCounters, got.distributedCounters)
	assert.Equal(test, getPointer(randomSource), getPointer(got.randomSource))
	assert.Equal(test, getPointer(formatter), getPointer(got.formatter))
}

func TestDistributedGenerator_GenerateCode(test *testing.T) {
	type fields struct {
		counter             counters.ChunkedCounter
		distributedCounters []counters.DistributedCounter
		randomSource        RandomSource
		formatter           Formatter
	}

	for _, data := range []struct {
		name        string
		fields      fields
		wantCounter counters.ChunkedCounter
		wantCode    string
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "success with incrementing",
			fields: fields{
				counter: func() counters.ChunkedCounter {
					counter := counters.NewChunkedCounter(23)
					counter.Reset(42)

					return counter
				}(),
				distributedCounters: []counters.DistributedCounter{
					NewMarkedDistributedCounter(1),
					NewMarkedDistributedCounter(2),
				},
				randomSource: func(maximum int) int { panic("not implemented") },
				formatter:    func(code uint64) string { return fmt.Sprintf("[%d]", code) },
			},
			wantCounter: func() counters.ChunkedCounter {
				counter := counters.NewChunkedCounter(23)
				counter.Reset(42)
				counter.Increase()

				return counter
			}(),
			wantCode: "[42]",
			wantErr:  assert.NoError,
		},
		{
			name: "success with resetting",
			fields: fields{
				counter: func() counters.ChunkedCounter {
					counter := counters.NewChunkedCounter(23)
					counter.Reset(42)
					for !counter.IsOver() {
						counter.Increase()
					}

					return counter
				}(),
				distributedCounters: func() []counters.DistributedCounter {
					firstCounter := NewMarkedDistributedCounter(1)

					secondCounter := NewMarkedDistributedCounter(2)
					secondCounter.On("NextCountChunk").Return(uint64(100), nil)

					return []counters.DistributedCounter{firstCounter, secondCounter}
				}(),
				randomSource: rand.New(rand.NewSource(1)).Intn,
				formatter:    func(code uint64) string { return fmt.Sprintf("[%d]", code) },
			},
			wantCounter: func() counters.ChunkedCounter {
				counter := counters.NewChunkedCounter(23)
				counter.Reset(100)
				counter.Increase()

				return counter
			}(),
			wantCode: "[100]",
			wantErr:  assert.NoError,
		},
		{
			name: "error with resetting",
			fields: fields{
				counter: func() counters.ChunkedCounter {
					counter := counters.NewChunkedCounter(23)
					counter.Reset(42)
					for !counter.IsOver() {
						counter.Increase()
					}

					return counter
				}(),
				distributedCounters: func() []counters.DistributedCounter {
					firstCounter := NewMarkedDistributedCounter(1)

					secondCounter := NewMarkedDistributedCounter(2)
					secondCounter.On("NextCountChunk").Return(uint64(0), iotest.ErrTimeout)

					return []counters.DistributedCounter{firstCounter, secondCounter}
				}(),
				randomSource: rand.New(rand.NewSource(1)).Intn,
				formatter:    func(code uint64) string { panic("not implemented") },
			},
			wantCounter: func() counters.ChunkedCounter {
				counter := counters.NewChunkedCounter(23)
				counter.Reset(42)
				for !counter.IsOver() {
					counter.Increase()
				}

				return counter
			}(),
			wantCode: "",
			wantErr:  assert.Error,
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
		counter             counters.ChunkedCounter
		distributedCounters []counters.DistributedCounter
		randomSource        RandomSource
		formatter           Formatter
	}

	for _, data := range []struct {
		name        string
		fields      fields
		wantCounter counters.ChunkedCounter
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				counter: func() counters.ChunkedCounter {
					counter := counters.NewChunkedCounter(23)
					counter.Reset(42)

					return counter
				}(),
				distributedCounters: func() []counters.DistributedCounter {
					firstCounter := NewMarkedDistributedCounter(1)

					secondCounter := NewMarkedDistributedCounter(2)
					secondCounter.On("NextCountChunk").Return(uint64(100), nil)

					return []counters.DistributedCounter{firstCounter, secondCounter}
				}(),
				randomSource: rand.New(rand.NewSource(1)).Intn,
				formatter:    func(code uint64) string { panic("not implemented") },
			},
			wantCounter: func() counters.ChunkedCounter {
				counter := counters.NewChunkedCounter(23)
				counter.Reset(100)

				return counter
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				counter: func() counters.ChunkedCounter {
					counter := counters.NewChunkedCounter(23)
					counter.Reset(42)

					return counter
				}(),
				distributedCounters: func() []counters.DistributedCounter {
					firstCounter := NewMarkedDistributedCounter(1)

					secondCounter := NewMarkedDistributedCounter(2)
					secondCounter.On("NextCountChunk").Return(uint64(0), iotest.ErrTimeout)

					return []counters.DistributedCounter{firstCounter, secondCounter}
				}(),
				randomSource: rand.New(rand.NewSource(1)).Intn,
				formatter:    func(code uint64) string { panic("not implemented") },
			},
			wantCounter: func() counters.ChunkedCounter {
				counter := counters.NewChunkedCounter(23)
				counter.Reset(42)

				return counter
			}(),
			wantErr: assert.Error,
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
	distributedCounters := []counters.DistributedCounter{
		NewMarkedDistributedCounter(1),
		NewMarkedDistributedCounter(2),
	}
	generator := &DistributedGenerator{
		counter:             counters.NewChunkedCounter(23),
		distributedCounters: distributedCounters,
		randomSource:        rand.New(rand.NewSource(1)).Intn,
		formatter:           func(code uint64) string { panic("not implemented") },
	}
	got := generator.selectCounter()

	for _, distributedCounter := range distributedCounters {
		mock.AssertExpectationsForObjects(test, distributedCounter)
	}
	assert.Equal(test, NewMarkedDistributedCounter(2), got)
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}
