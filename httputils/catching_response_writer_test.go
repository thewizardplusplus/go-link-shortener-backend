package httputils

import (
	"errors"
	"net/http"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCatchingResponseWriter(test *testing.T) {
	writer := new(MockResponseWriter)
	got := NewCatchingResponseWriter(writer)

	mock.AssertExpectationsForObjects(test, writer)
	assert.Equal(test, writer, got.ResponseWriter)
	assert.NoError(test, got.lastError)
}

func TestCatchingResponseWriter_LastError(test *testing.T) {
	type fields struct {
		lastError error
	}

	for _, data := range []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "without the last error",
			fields: fields{
				lastError: nil,
			},
			wantErr: nil,
		},
		{
			name: "with the last error",
			fields: fields{
				lastError: iotest.ErrTimeout,
			},
			wantErr: iotest.ErrTimeout,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			writer := CatchingResponseWriter{
				lastError: data.fields.lastError,
			}
			gotErr := writer.LastError()

			assert.Equal(test, data.wantErr, gotErr)
		})
	}
}

func TestCatchingResponseWriter_Write(test *testing.T) {
	type fields struct {
		responseWriter http.ResponseWriter
		lastError      error
	}
	type args struct {
		p []byte
	}

	for _, data := range []struct {
		name        string
		fields      fields
		args        args
		wantN       int
		wantErr     assert.ErrorAssertionFunc
		wantLastErr error
	}{
		{
			name: "success",
			fields: fields{
				responseWriter: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Write", []byte("test")).Return(4, nil)

					return writer
				}(),
				lastError: nil,
			},
			args:        args{[]byte("test")},
			wantN:       4,
			wantErr:     assert.NoError,
			wantLastErr: nil,
		},
		{
			name: "error (without the last error)",
			fields: fields{
				responseWriter: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Write", []byte("test")).Return(2, iotest.ErrTimeout)

					return writer
				}(),
				lastError: nil,
			},
			args:        args{[]byte("test")},
			wantN:       2,
			wantErr:     assert.Error,
			wantLastErr: iotest.ErrTimeout,
		},
		{
			name: "error (with the last error)",
			fields: fields{
				responseWriter: func() http.ResponseWriter {
					writer := new(MockResponseWriter)
					writer.On("Write", []byte("test")).Return(2, iotest.ErrTimeout)

					return writer
				}(),
				lastError: errors.New("dummy"),
			},
			args:        args{[]byte("test")},
			wantN:       2,
			wantErr:     assert.Error,
			wantLastErr: iotest.ErrTimeout,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			writer := &CatchingResponseWriter{
				ResponseWriter: data.fields.responseWriter,
				lastError:      data.fields.lastError,
			}
			gotN, gotErr := writer.Write(data.args.p)

			mock.AssertExpectationsForObjects(test, data.fields.responseWriter)
			assert.Equal(test, data.wantN, gotN)
			data.wantErr(test, gotErr)
			assert.Equal(test, data.wantLastErr, writer.lastError)
		})
	}
}
