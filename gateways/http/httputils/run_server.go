package httputils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/go-log/log"
)

//go:generate mockery -name=Server -inpkg -case=underscore -testonly

// Server ...
type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// RunServerDependencies ...
type RunServerDependencies struct {
	Server Server
	Logger log.Logger
}

// RunServer ...
func RunServer(
	shutdownCtx context.Context,
	dependencies RunServerDependencies,
	interruptSignals ...os.Signal,
) (ok bool) {
	var waiter sync.WaitGroup
	waiter.Add(1)

	interruptCtx, interruptCtxCancel := context.WithCancel(context.Background())
	defer interruptCtxCancel()

	go func() {
		defer waiter.Done()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, interruptSignals...)

		select {
		case <-interrupt:
		case <-interruptCtx.Done():
			return
		}

		err := dependencies.Server.Shutdown(shutdownCtx)
		if ok = err == nil; !ok {
			// error with closing listeners
			dependencies.Logger.Logf("error with shutdown: %v", err)
		}
	}()

	if err := dependencies.Server.ListenAndServe(); err != http.ErrServerClosed {
		// error with starting or closing listeners
		dependencies.Logger.Logf("error with listening and serving: %v", err)
		return false
	}

	waiter.Wait()
	return ok
}
