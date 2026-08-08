package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"scaling-up-rest-vs-grpc/internal/data/cache"
	"scaling-up-rest-vs-grpc/internal/lib/redis"
	"scaling-up-rest-vs-grpc/internal/rest"
)

const (
	http1Addr = ":8080"
	http2Addr = ":8081"
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

	cached := cache.New()
	if err := cached.LoadFromRedis(context.Background(), client); err != nil {
		slog.Error("failed to load dataset from redis", "error", err)
		os.Exit(1)
	}
	if cached.GetSmallDataset() == nil || cached.GetLargeDataset() == nil {
		slog.Warn("dataset not seeded yet")
	}

	shapeCached := cache.NewShape()
	if err := shapeCached.LoadFromRedis(context.Background(), client); err != nil {
		slog.Error("failed to load shape experiment dataset from redis", "error", err)
		os.Exit(1)
	}

	http1Server, http2Server := rest.NewServers(http1Addr, http2Addr, cached, shapeCached)

	errCh := make(chan error, 2)

	go func() {
		slog.Info("rest-server (HTTP/1.1) started", "addr", http1Addr)
		errCh <- http1Server.ListenAndServe()
	}()

	go func() {
		slog.Info("rest-server (HTTP/2 cleartext) started", "addr", http2Addr)
		errCh <- http2Server.ListenAndServe()
	}()

	if err := <-errCh; err != nil {
		slog.Error("rest server stopped", "error", err)
		os.Exit(1)
	}
}
