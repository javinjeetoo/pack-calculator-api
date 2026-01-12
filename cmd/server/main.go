package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/javinjeetoo/pack-calculator-api/internal/api"
	"github.com/javinjeetoo/pack-calculator-api/internal/ui"
)

func main() {
	// Default pack sizes
	defaultPackSizes := []int{250, 500, 1000, 2000, 5000}
	if env := os.Getenv("PACK_SIZES"); env != "" {
		if parsed, err := parsePackSizes(env); err == nil && len(parsed) > 0 {
			defaultPackSizes = parsed
		} else if err != nil {
			log.Printf("WARN: invalid PACK_SIZES, using defaults: %v", err)
		}
	}

	// Port configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	h := api.NewHandler(defaultPackSizes)

	// HTTP routing
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/calculate", h.Calculate)

	// UI served at /
	uiHandler, err := ui.Handler()
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", uiHandler)

	// HTTP server with sensible timeouts
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Context that listens for SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on :%s", port)
		errCh <- server.ListenAndServe()
	}()

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}

	// Attempt graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		_ = server.Close()
	}

	log.Println("server stopped")
}

func parsePackSizes(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}
