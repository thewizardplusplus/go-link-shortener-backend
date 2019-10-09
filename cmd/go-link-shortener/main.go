package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	"github.com/thewizardplusplus/go-link-shortener/code"
	"github.com/thewizardplusplus/go-link-shortener/entities"
	"github.com/thewizardplusplus/go-link-shortener/gateways/cache"
	"github.com/thewizardplusplus/go-link-shortener/gateways/handlers"
	"github.com/thewizardplusplus/go-link-shortener/gateways/presenters"
	"github.com/thewizardplusplus/go-link-shortener/gateways/router"
	"github.com/thewizardplusplus/go-link-shortener/gateways/storage"
	"github.com/thewizardplusplus/go-link-shortener/usecases"
)

// nolint: lll
type options struct {
	ServerAddress  string `env:"SERVER_ADDRESS" envDefault:":8080"`
	CacheAddress   string `env:"CACHE_ADDRESS" envDefault:"localhost:6379"`
	StorageAddress string `env:"STORAGE_ADDRESS" envDefault:"mongodb://localhost:27017"`
}

const (
	storageDatabase   = "go-link-shortener"
	storageCollection = "links"
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
					CodeGenerator: code.UUIDGenerator{},
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
