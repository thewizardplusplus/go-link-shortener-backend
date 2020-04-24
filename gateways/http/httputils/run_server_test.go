package httputils

import (
	"context"
	"net/http"
	"os"
	"testing"
	"testing/iotest"
	"time"

	"github.com/go-log/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunServer(test *testing.T) {
	type args struct {
		shutdownCtx      context.Context
		dependencies     RunServerDependencies
		interruptSignals []os.Signal
	}

	for _, data := range []struct {
		name   string
		args   args
		action func(test *testing.T)
		wantOk assert.BoolAssertionFunc
	}{
		{
			name: "success",
			args: args{
				shutdownCtx: context.Background(),
				dependencies: RunServerDependencies{
					Server: func() Server {
						server := new(MockServer)
						server.On("ListenAndServe").Return(http.ErrServerClosed)
						server.
							On(
								"Shutdown",
								mock.MatchedBy(func(context.Context) bool { return true }),
							).
							Return(nil)

						return server
					}(),
					Logger: new(MockLogger),
				},
				interruptSignals: []os.Signal{os.Interrupt},
			},
			action: func(test *testing.T) {
				time.Sleep(time.Second)

				currentProcess, err := os.FindProcess(os.Getpid())
				require.NoError(test, err)

				err = currentProcess.Signal(os.Interrupt)
				require.NoError(test, err)
			},
			wantOk: assert.True,
		},
		{
			name: "error on the ListenAndServe() call",
			args: args{
				shutdownCtx: context.Background(),
				dependencies: RunServerDependencies{
					Server: func() Server {
						server := new(MockServer)
						server.On("ListenAndServe").Return(iotest.ErrTimeout)

						return server
					}(),
					Logger: func() log.Logger {
						logger := new(MockLogger)
						logger.
							On(
								"Logf",
								mock.MatchedBy(func(string) bool { return true }),
								iotest.ErrTimeout,
							).
							Return()

						return logger
					}(),
				},
				interruptSignals: []os.Signal{os.Interrupt},
			},
			action: func(test *testing.T) {},
			wantOk: assert.False,
		},
		{
			name: "error on the Shutdown() call",
			args: args{
				shutdownCtx: context.Background(),
				dependencies: RunServerDependencies{
					Server: func() Server {
						server := new(MockServer)
						server.On("ListenAndServe").Return(http.ErrServerClosed)
						server.
							On(
								"Shutdown",
								mock.MatchedBy(func(context.Context) bool { return true }),
							).
							Return(iotest.ErrTimeout)

						return server
					}(),
					Logger: func() log.Logger {
						logger := new(MockLogger)
						logger.
							On(
								"Logf",
								mock.MatchedBy(func(string) bool { return true }),
								iotest.ErrTimeout,
							).
							Return()

						return logger
					}(),
				},
				interruptSignals: []os.Signal{os.Interrupt},
			},
			action: func(test *testing.T) {
				time.Sleep(time.Second)

				currentProcess, err := os.FindProcess(os.Getpid())
				require.NoError(test, err)

				err = currentProcess.Signal(os.Interrupt)
				require.NoError(test, err)
			},
			wantOk: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			go data.action(test)

			gotOk := RunServer(
				data.args.shutdownCtx,
				data.args.dependencies,
				data.args.interruptSignals...,
			)

			mock.AssertExpectationsForObjects(
				test,
				data.args.dependencies.Server,
				data.args.dependencies.Logger,
			)
			data.wantOk(test, gotOk)
		})
	}
}
