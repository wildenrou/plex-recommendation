package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/config"
	httpinternal "github.com/wgeorgecook/plex-recommendation/internal/pkg/http"
)

func main() {
	log.Println("Hello!")
	defer log.Println("Good bye!")

	ctx := context.Background()

	// Set up open telemetry
	log.Println("initializing open telemetry client...")
	shutdownOtel, err := telemetry.InitOtel(ctx,
		telemetry.WithTracer(true),
		telemetry.WithMeter(false))
	if err != nil {
		panic(err)
	}
	log.Println("done!")
	defer func() {
		if err := shutdownOtel(ctx); err != nil {
			log.Println("could not shutdown otel:" + err.Error())
		}
	}()

	// Start the server in a goroutine
	serverDone := make(chan error, 1)
	go func() {
		log.Println("Starting server...")
		httpinternal.StartServer(ctx, config.LoadConfig(), serverDone)
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
