package main

import (
	"context"
	"time"
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

}
