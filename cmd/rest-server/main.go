package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"scaling-up-rest-vs-grpc/internal/data/store"
	"scaling-up-rest-vs-grpc/internal/lib/redis"
	"scaling-up-rest-vs-grpc/internal/rest"
)

const (
	listenAddr = ":8443"
	certPath   = "certs/server.crt"
	keyPath    = "certs/server.key"
)

func main() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		slog.Error("REDIS_ADDR is not set")
		os.Exit(1)
	}

	client := redis.NewClient(addr)
	if err := client.WaitReady(context.Background(), 10, 2*time.Second); err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	dataStore := store.New()
	if err := dataStore.LoadFromRedis(context.Background(), client); err != nil {
		slog.Error("failed to load dataset from redis", "error", err)
		os.Exit(1)
	}
	if dataStore.GetSmallDataset() == nil || dataStore.GetLargeDataset() == nil {
		slog.Warn("dataset not seeded yet")
	}

	server := rest.NewServer(listenAddr, dataStore)
	slog.Info("rest-server started", "addr", listenAddr)

	if err := server.ListenAndServeTLS(certPath, keyPath); err != nil {
		slog.Error("rest server stopped", "error", err)
		os.Exit(1)
	}
}
