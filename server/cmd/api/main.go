package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"modmapper/server/internal/handlers"
)

func main() {
	// init config + DB client
	cfg := MustLoadConfig()
	mongoClient, db := MustConnectMongo(cfg)
	defer func() { // defer this anon fn
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = mongoClient.Disconnect(ctx)
	}()

	// build router
	router := handlers.NewRouter(handlers.Config{
		CORSOrigin: cfg.CORSOrigin,
	}, db)

	// server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// start server
	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}

	}()

	stop := make(chan os.Signal, 1)                      // chan is a channel type constructor (built-in keyword). The channel can carry the data type that comes after
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM) // notify the stop channel when signals listed here are received
	<-stop                                               // line blocks until the correct signal is received

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown errer: %v", err)
	}
}
