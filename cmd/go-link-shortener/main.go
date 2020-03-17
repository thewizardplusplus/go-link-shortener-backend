package main

// nolint: lll
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
	"github.com/thewizardplusplus/go-link-shortener-backend/entities"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/cache"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/counter"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/http/handlers"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/http/handlers/presenters"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/http/httputils"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/http/router"
	"github.com/thewizardplusplus/go-link-shortener-backend/gateways/storage"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators/counters"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators/counters/transformers"
	"github.com/thewizardplusplus/go-link-shortener-backend/usecases/generators/formatters"
)

type options struct {
	Server struct {
		Address    string `env:"SERVER_ADDRESS" envDefault:":8080"`
		StaticPath string `env:"SERVER_STATIC_PATH" envDefault:"./static"`
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
		Range   uint64 `env:"COUNTER_RANGE" envDefault:"1000000000"`
	}
}

const (
	errorURL               = "/error"
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

	var distributedCounters []counters.DistributedCounter
	for i := 0; i < options.Counter.Count; i++ {
		distributedCounters = append(distributedCounters, counters.TransformedCounter{
			DistributedCounter: counter.Counter{
				Client: counterClient,
				Name:   fmt.Sprintf(counterNameTemplate, i),
			},
			Transformer: transformers.NewLinear(
				transformers.WithFactor(options.Counter.Chunk),
				transformers.WithOffset(uint64(i)*options.Counter.Range),
			),
		})
	}

	linkByCodeGetter := usecases.LinkGetterGroup{
		cacheGetter,
		storage.LinkGetter{
			Client:     storageClient,
			Database:   storageDatabase,
			Collection: storageCollection,
			KeyField:   "code",
		},
	}

	redirectPresenter := presenters.RedirectPresenter{
		ErrorURL: errorURL,
		Printer:  errorLogger,
	}
	jsonLinkPresenter := presenters.SilentLinkPresenter{
		LinkPresenter: presenters.JSONPresenter{},
		Printer:       errorLogger,
	}
	jsonErrorPresenter := presenters.SilentErrorPresenter{
		ErrorPresenter: presenters.JSONPresenter{},
		Printer:        errorLogger,
	}

	routerHandler := router.NewRouter(redirectEndpointPrefix, router.Handlers{
		LinkRedirectHandler: handlers.LinkGettingHandler{
			LinkGetter: linkByCodeGetter,
			LinkPresenter: presenters.SilentLinkPresenter{
				LinkPresenter: redirectPresenter,
				Printer:       errorLogger,
			},
			ErrorPresenter: presenters.SilentErrorPresenter{
				ErrorPresenter: redirectPresenter,
				Printer:        errorLogger,
			},
		},
		LinkGettingHandler: handlers.LinkGettingHandler{
			LinkGetter:     linkByCodeGetter,
			LinkPresenter:  jsonLinkPresenter,
			ErrorPresenter: jsonErrorPresenter,
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
				CodeGenerator: generators.NewDistributedGenerator(
					options.Counter.Chunk,
					counters.CounterGroup{
						DistributedCounters: distributedCounters,
						RandomSource:        rand.New(rand.NewSource(time.Now().UnixNano())).Intn,
					},
					formatters.InBase62,
				),
			},
			LinkPresenter:  jsonLinkPresenter,
			ErrorPresenter: jsonErrorPresenter,
		},
		StaticFileHandler: httputils.CatchingHandler{
			Handler: http.FileServer(http.Dir(options.Server.StaticPath)),
			Printer: errorLogger,
		},
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
