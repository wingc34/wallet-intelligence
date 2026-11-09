package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status": "ok",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode health response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	stop, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	serverErr := make(chan error, 1)

	// ListenAndServe() is a blocking operation, use goroutine for handling, except blocking
	go func() {
		log.Println("API server starting on :8080")
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	case <-stop.Done():
		log.Println("Shutting down HTTP server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)

			if closeErr := server.Close(); closeErr != nil {
				log.Printf("Failed to close HTTP server: %v", closeErr)
			}
		}

		if err := <-serverErr; !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}

		log.Printf("HTTP server stopped")

	}
}
