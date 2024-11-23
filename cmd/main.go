package main

import (
	"context"
	"errors"
	"github.com/ntsiris/todo-app/internal/store"
	"github.com/ntsiris/todo-app/internal/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	storage := store.NewMapStore()

	err := storage.Open()

	if err != nil {
		log.Fatal(err)
	}
	defer func(storage *store.MapStore) {
		err := storage.Close()
		if err != nil {

		}
	}(storage)

	err = storage.VerifyConnection()
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	apiServerAddr := ":8080"
	apiServer := transport.NewAPIServer(apiServerAddr, storage)

	go func() {
		log.Println("Listening on ", apiServerAddr)
		if err := apiServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("api server run failed: %v", err)
		}
	}()

	// Wait for application to terminate
	<-stop
	log.Println("Shutting down server...")

	// Create context with timeout for shutting down process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}

	// Perform any cleanup necessary
}
