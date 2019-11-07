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

type counterOptions struct {
	Address string `env:"COUNTER_ADDRESS" envDefault:"localhost:2379"`
	Count   int    `env:"COUNTER_COUNT" envDefault:"2"`
	Chunk   uint64 `env:"COUNTER_CHUNK" envDefault:"1000"`
}

// nolint: lll
type options struct {
	ServerAddress  string `env:"SERVER_ADDRESS" envDefault:":8080"`
	CacheAddress   string `env:"CACHE_ADDRESS" envDefault:"localhost:6379"`
	StorageAddress string `env:"STORAGE_ADDRESS" envDefault:"mongodb://localhost:27017"`
	Counter        counterOptions
}

const (
	storageDatabase     = "go-link-shortener"
	storageCollection   = "links"
	counterNameTemplate = "distributed-counter-%d"
)

func main() {
	var options options // nolint: vetshadow
	if err := env.Parse(&options); err != nil {
		log.Fatalf("error on parsing options: %v", err)
	}

	cacheClient := cache.NewClient(options.CacheAddress)
	cacheGetter := cache.LinkGetter{Client: cacheClient}

	storageClient, err := storage.NewClient(options.StorageAddress)
	if err != nil {
		log.Fatalf("error on creating the storage client: %v", err)
	}

	counterClient, err := counter.NewClient(options.Counter.Address)
	if err != nil {
		log.Fatalf("error on creating the counter client: %v", err)
	}

	var counters []code.DistributedCounter
	for i := 0; i < options.Counter.Count; i++ {
		counters = append(counters, counter.Counter{
			Client: counterClient,
			Name:   fmt.Sprintf(counterNameTemplate, i),
		})
	}

	var presenter presenters.JSONPresenter
	server := http.Server{
		Addr: options.ServerAddress,
		Handler: router.NewRouter(router.Handlers{
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
				LinkPresenter:  presenter,
				ErrorPresenter: presenter,
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
						cache.LinkSetter{
							KeyExtractor: func(link entities.Link) string { return link.Code },
							Client:       cacheClient,
							Expiration:   time.Hour,
						},
						cache.LinkSetter{
							KeyExtractor: func(link entities.Link) string { return link.URL },
							Client:       cacheClient,
							Expiration:   time.Hour,
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
				LinkPresenter:  presenter,
				ErrorPresenter: presenter,
			},
			NotFoundHandler: handlers.NotFoundHandler{ErrorPresenter: presenter},
		}),
	}

	done := make(chan struct{})
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		<-interrupt

		if err := server.Shutdown(context.Background()); err != nil {
			// error on closing listeners
			log.Printf("error on shutdown: %v", err)
		}

		close(done)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// error on starting or closing listeners
		log.Fatalf("error on listening and serving: %v", err)
	}

	<-done
}
