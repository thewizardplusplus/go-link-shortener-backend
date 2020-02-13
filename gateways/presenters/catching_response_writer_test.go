package presenters

import (
	"net/http"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_newCatchingResponseWriter(test *testing.T) {
	writer := new(MockResponseWriter)
	got := newCatchingResponseWriter(writer)

	mock.AssertExpectationsForObjects(test, writer)
	assert.Equal(test, writer, got.ResponseWriter)
	assert.NoError(test, got.error)
}

func Test_catchingResponseWriter_Write(test *testing.T) {
	type fields struct {
		responseWriter http.ResponseWriter
	}
	type args struct {
		p []byte
	}

	for _, data := range []struct {
		name           string
		fields         fields
		args           args
		wantN          int
		wantErr        assert.ErrorAssertionFunc
		wantCatchedErr error
	}{
		{
			name: "success",
			fields: fields{
				responseWriter: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Write", []byte("test")).Return(4, nil)

					return writer
				}(),
			},
			args:           args{[]byte("test")},
			wantN:          4,
			wantErr:        assert.NoError,
			wantCatchedErr: nil,
		},
		{
			name: "error",
			fields: fields{
				responseWriter: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Write", []byte("test")).Return(2, iotest.ErrTimeout)

					return writer
				}(),
			},
			args:           args{[]byte("test")},
			wantN:          2,
			wantErr:        assert.Error,
			wantCatchedErr: iotest.ErrTimeout,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			writer := &catchingResponseWriter{
				ResponseWriter: data.fields.responseWriter,
			}
			gotN, gotErr := writer.Write(data.args.p)

			mock.AssertExpectationsForObjects(test, data.fields.responseWriter)
			assert.Equal(test, data.wantN, gotN)
			data.wantErr(test, gotErr)
			assert.Equal(test, data.wantCatchedErr, writer.error)
		})
	}
}
