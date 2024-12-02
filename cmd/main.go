package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/ntsiris/todo-app/internal/store"
	"github.com/ntsiris/todo-app/internal/transport"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, using environment variables")
	}

	var storeConfig store.StoreConfig
	if err := envconfig.Process("", &storeConfig); err != nil {
		log.Fatalf("Failed to process storage configuration from environment: %v", err)
	}

	storage, err := store.NewPostgresStore(&storeConfig)
	if err != nil {
		log.Fatalf("Failed to create postgres storage: %v", err)
	}

	err = storage.Open()
	if err != nil {
		log.Fatalf("Failed to open postgres storage: %v", err)
	}
	defer func() {
		err := storage.Close()
		if err != nil {

		}
	}()

	err = storage.VerifyConnection()
	if err != nil {
		log.Fatalf("Failed to verify connection to postgres storage: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var apiServerConfig transport.APIServerConfig
	if err := envconfig.Process("", &apiServerConfig); err != nil {
		log.Fatalf("Failed to process api server configuration: %v", err)
	}

	apiServer := transport.NewAPIServer(&apiServerConfig, storage)

	go func() {
		log.Println("Listening on ", fmt.Sprintf("%s:%d", apiServerConfig.Host, apiServerConfig.Port))
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
	if err := apiServer.Stop(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}

	// Perform any cleanup necessary
}
