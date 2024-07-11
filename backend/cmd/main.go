package main

import (
	"context"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/config"
	httpinternal "github.com/wgeorgecook/plex-recommendation/internal/pkg/http"
)

func main() {
	log.Println("Hello!")
	defer log.Println("Good bye!")

	// Set up open telemetry
	log.Println("initializing open telemetry client...")
	shutdownOtel, err := telemetry.InitOtel(context.Background(), telemetry.WithTracer(true))
	if err != nil {
		panic(err)
	}
	log.Println("done!")
	defer func() {
		if err := shutdownOtel(context.Background()); err != nil {
			log.Println("could not shutdown otel:" + err.Error())
		}
	}()

	// Start the server in a goroutine
	serverDone := make(chan error, 1)
	go func() {
		log.Println("Starting server...")
		httpinternal.StartServer(config.LoadConfig(), serverDone)
	}()

	// Set up signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal or the server exits
	select {
	case err := <-serverDone:
		if err != nil {
			log.Fatalf("Server error: %v", err)
		} else {
			log.Println("Server exited gracefully")
		}
	case <-quit:
		log.Println("Shutdown signal received, initiating graceful shutdown...")
		log.Println("Server shutdown complete")
	}
}
