package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	serverAddress = ":8080"
)

func main() {
	server := http.Server{
		Addr: serverAddress,
		Handler: http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			io.WriteString(writer, "Hello, world!\n")
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
