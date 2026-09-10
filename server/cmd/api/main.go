package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"modmapper/server/internal/platform/config"
	"modmapper/server/internal/platform/database"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load("../.env")
	if err != nil {
		return err
	}

	startupCtx, cancelStartup := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	pool, err := database.Open(startupCtx, cfg.DatabaseURL)
	cancelStartup()

	if err != nil {
		return err
	}
	defer pool.Close()

	mux := http.NewServeMux()
	// health check to see if http server is running
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// report whether db is currently reachable
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			http.Error(
				w,
				"database unavailable",
				http.StatusServiceUnavailable,
			)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ready",
		})
	})

	server := &http.Server{
		Addr:              net.JoinHostPort("127.0.0.1", cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second, // keep-alive conn is closed after this
	}

	// listen for any termination signal
	stopCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("API listening at http://%s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err

	case <-stopCtx.Done():
		log.Print("Shutting down API")
		shutdownCtx, cancelShutdown := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancelShutdown()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close() // force close instead of graceful shutdown
			return err
		}

		return nil
	}
}
