package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	middlewares "github.com/gorilla/handlers"
	"github.com/thewizardplusplus/go-link-shortener/code"
	"github.com/thewizardplusplus/go-link-shortener/entities"
	"github.com/thewizardplusplus/go-link-shortener/gateways/cache"
	"github.com/thewizardplusplus/go-link-shortener/gateways/counter"
	"github.com/thewizardplusplus/go-link-shortener/gateways/handlers"
	"github.com/thewizardplusplus/go-link-shortener/gateways/presenters"
	"github.com/thewizardplusplus/go-link-shortener/gateways/router"
	"github.com/thewizardplusplus/go-link-shortener/gateways/storage"
	"github.com/thewizardplusplus/go-link-shortener/usecases"
)

type options struct {
	Server struct {
		Address string `env:"SERVER_ADDRESS" envDefault:":8080"`
	}
	Cache struct {
		Address string `env:"CACHE_ADDRESS" envDefault:"localhost:6379"`
		TTL     struct {
			Code time.Duration `env:"CACHE_TTL_CODE" envDefault:"1h"`
			URL  time.Duration `env:"CACHE_TTL_URL" envDefault:"1h"`
		}
	}
	Storage struct {
		Address string `env:"STORAGE_ADDRESS" envDefault:"mongodb://localhost:27017"`
	}
	Counter struct {
		Address string `env:"COUNTER_ADDRESS" envDefault:"localhost:2379"`
		Count   int    `env:"COUNTER_COUNT" envDefault:"2"`
		Chunk   uint64 `env:"COUNTER_CHUNK" envDefault:"1000"`
	}
}

const (
	redirectEndpointPrefix = "/redirect"
	storageDatabase        = "go-link-shortener"
	storageCollection      = "links"
	counterNameTemplate    = "distributed-counter-%d"
)

func main() {
	errorLogger := log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)

	var options options // nolint: vetshadow
	if err := env.Parse(&options); err != nil {
		errorLogger.Fatalf("error with parsing options: %v", err)
	}

	cacheClient := cache.NewClient(options.Cache.Address)
	cacheGetter := usecases.SilentLinkGetter{
		LinkGetter: cache.LinkGetter{Client: cacheClient},
		Printer:    errorLogger,
	}

	storageClient, err := storage.NewClient(options.Storage.Address)
	if err != nil {
		errorLogger.Fatalf("error with creating the storage client: %v", err)
	}

	counterClient, err := counter.NewClient(options.Counter.Address)
	if err != nil {
		errorLogger.Fatalf("error with creating the counter client: %v", err)
	}

	var counters []code.DistributedCounter
	for i := 0; i < options.Counter.Count; i++ {
		counters = append(counters, counter.Counter{
			Client: counterClient,
			Name:   fmt.Sprintf(counterNameTemplate, i),
		})
	}

	linkPresenter := presenters.SilentLinkPresenter{
		LinkPresenter: presenters.JSONPresenter{},
		Printer:       errorLogger,
	}
	errorPresenter := presenters.SilentErrorPresenter{
		ErrorPresenter: presenters.JSONPresenter{},
		Printer:        errorLogger,
	}

	routerHandler := router.NewRouter(redirectEndpointPrefix, router.Handlers{
		LinkGettingHandler: handlers.LinkGettingHandler{
			LinkGetter: usecases.LinkGetterGroup{
				cacheGetter,
				storage.LinkGetter{
					Client:     storageClient,
					Database:   storageDatabase,
					Collection: storageCollection,
					KeyField:   "code",
				},
			},
			LinkPresenter:  linkPresenter,
			ErrorPresenter: errorPresenter,
		},
		LinkCreatingHandler: handlers.LinkCreatingHandler{
			LinkCreator: usecases.LinkCreator{
				LinkGetter: usecases.LinkGetterGroup{
					cacheGetter,
					storage.LinkGetter{
						Client:     storageClient,
						Database:   storageDatabase,
						Collection: storageCollection,
						KeyField:   "url",
					},
				},
				LinkSetter: usecases.LinkSetterGroup{
					usecases.SilentLinkSetter{
						LinkSetter: cache.LinkSetter{
							KeyExtractor: func(link entities.Link) string { return link.Code },
							Client:       cacheClient,
							Expiration:   options.Cache.TTL.Code,
						},
						Printer: errorLogger,
					},
					usecases.SilentLinkSetter{
						LinkSetter: cache.LinkSetter{
							KeyExtractor: func(link entities.Link) string { return link.URL },
							Client:       cacheClient,
							Expiration:   options.Cache.TTL.URL,
						},
						Printer: errorLogger,
					},
					storage.LinkSetter{
						Client:     storageClient,
						Database:   storageDatabase,
						Collection: storageCollection,
					},
				},
				CodeGenerator: code.NewDistributedGenerator(
					options.Counter.Chunk,
					counters,
					rand.New(rand.NewSource(time.Now().UnixNano())).Intn,
				),
			},
			LinkPresenter:  linkPresenter,
			ErrorPresenter: errorPresenter,
		},
		NotFoundHandler: handlers.NotFoundHandler{ErrorPresenter: errorPresenter},
	})
	routerHandler.
		Use(middlewares.RecoveryHandler(middlewares.RecoveryLogger(errorLogger)))
	routerHandler.
		Use(func(handler http.Handler) http.Handler {
			return middlewares.LoggingHandler(os.Stdout, handler)
		})

	server := http.Server{
		Addr:    options.Server.Address,
		Handler: routerHandler,
	}

	done := make(chan struct{})
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		<-interrupt

		if err := server.Shutdown(context.Background()); err != nil {
			// error with closing listeners
			errorLogger.Printf("error with shutdown: %v", err)
		}

		close(done)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// error with starting or closing listeners
		errorLogger.Fatalf("error with listening and serving: %v", err)
	}

	<-done
}
