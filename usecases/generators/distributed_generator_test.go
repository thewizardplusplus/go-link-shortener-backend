package generators

// nolint: lll
import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators/counters"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators/formatters"
)

type MemorableDistributedCounter struct {
	MockDistributedCounterGroup

	Counters []counters.DistributedCounter
}

func (
	group *MemorableDistributedCounter,
) SelectCounter() counters.DistributedCounter {
	counter := group.MockDistributedCounterGroup.SelectCounter()
	group.Counters = append(group.Counters, counter)

	return counter
}

func TestNewDistributedGenerator(test *testing.T) {
	distributedCounters := new(MemorableDistributedCounter)
	formatter := func(code uint64) string { panic("not implemented") }
	got := NewDistributedGenerator(23, distributedCounters, formatter)

	mock.AssertExpectationsForObjects(test, distributedCounters)
	require.NotNil(test, got)
	assert.Equal(test, counters.NewChunkedCounter(23), got.counter)
	assert.Equal(test, distributedCounters, got.distributedCounters)
	assert.Equal(test, getPointer(formatter), getPointer(got.formatter))
}

func TestDistributedGenerator_GenerateCode(test *testing.T) {
	type fields struct {
		counter             counters.ChunkedCounter
		distributedCounters DistributedCounterGroup
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
				distributedCounters: new(MemorableDistributedCounter),
				formatter: func(code uint64) string {
					return fmt.Sprintf("[%d]", code)
				},
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
				distributedCounters: func() DistributedCounterGroup {
					counter := new(MockDistributedCounter)
					counter.On("NextCountChunk").Return(uint64(100), nil)

					group := new(MemorableDistributedCounter)
					group.On("SelectCounter").Return(counter)

					return group
				}(),
				formatter: func(code uint64) string { return fmt.Sprintf("[%d]", code) },
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
				distributedCounters: func() DistributedCounterGroup {
					counter := new(MockDistributedCounter)
					counter.On("NextCountChunk").Return(uint64(0), iotest.ErrTimeout)

					group := new(MemorableDistributedCounter)
					group.On("SelectCounter").Return(counter)

					return group
				}(),
				formatter: func(code uint64) string { panic("not implemented") },
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
				formatter:           data.fields.formatter,
			}
			gotCode, gotErr := generator.GenerateCode()

			mock.AssertExpectationsForObjects(test, data.fields.distributedCounters)
			counters :=
				data.fields.distributedCounters.(*MemorableDistributedCounter).Counters
			for _, counter := range counters {
				mock.AssertExpectationsForObjects(test, counter)
			}
			assert.Equal(test, data.wantCounter, generator.counter)
			assert.Equal(test, data.wantCode, gotCode)
			data.wantErr(test, gotErr)
		})
	}
}

func TestDistributedGenerator_GenerateCode_bulky(test *testing.T) {
	var countChunkOne uint64
	counterOne := new(MockDistributedCounter)
	counterOne.
		On("NextCountChunk").
		Return(
			func() uint64 {
				defer func() { countChunkOne++ }()
				return countChunkOne * 10
			},
			nil,
		)

	var countChunkTwo uint64
	counterTwo := new(MockDistributedCounter)
	counterTwo.
		On("NextCountChunk").
		Return(
			func() uint64 {
				defer func() { countChunkTwo++ }()
				return countChunkTwo*10 + 100
			},
			nil,
		)

	generator := NewDistributedGenerator(
		10,
		counters.CounterGroup{
			DistributedCounters: []counters.DistributedCounter{counterOne, counterTwo},
			RandomSource:        rand.New(rand.NewSource(1)).Intn,
		},
		formatters.InBase10,
	)

	var gotCodes []string
	for i := 0; i < 100; i++ {
		code, err := generator.GenerateCode()
		require.NoError(test, err)

		gotCodes = append(gotCodes, code)
	}

	var wantCodes []string
	var counterIndex, counterValueOne, counterValueTwo int
	random := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		if i%10 == 0 {
			counterIndex = random.Intn(2)
		}

		var code int
		switch counterIndex {
		case 0:
			code = counterValueOne
			counterValueOne++
		case 1:
			code = counterValueTwo + 100
			counterValueTwo++
		}

		wantCodes = append(wantCodes, strconv.Itoa(code))
	}

	assert.Equal(test, wantCodes, gotCodes)
}

func TestDistributedGenerator_resetCounter(test *testing.T) {
	type fields struct {
		counter             counters.ChunkedCounter
		distributedCounters DistributedCounterGroup
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
				distributedCounters: func() DistributedCounterGroup {
					counter := new(MockDistributedCounter)
					counter.On("NextCountChunk").Return(uint64(100), nil)

					group := new(MemorableDistributedCounter)
					group.On("SelectCounter").Return(counter)

					return group
				}(),
				formatter: func(code uint64) string { panic("not implemented") },
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
				distributedCounters: func() DistributedCounterGroup {
					counter := new(MockDistributedCounter)
					counter.On("NextCountChunk").Return(uint64(0), iotest.ErrTimeout)

					group := new(MemorableDistributedCounter)
					group.On("SelectCounter").Return(counter)

					return group
				}(),
				formatter: func(code uint64) string { panic("not implemented") },
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
				formatter:           data.fields.formatter,
			}
			gotErr := generator.resetCounter()

			mock.AssertExpectationsForObjects(test, data.fields.distributedCounters)
			counters :=
				data.fields.distributedCounters.(*MemorableDistributedCounter).Counters
			for _, counter := range counters {
				mock.AssertExpectationsForObjects(test, counter)
			}
			assert.Equal(test, data.wantCounter, generator.counter)
			data.wantErr(test, gotErr)
		})
	}
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}
