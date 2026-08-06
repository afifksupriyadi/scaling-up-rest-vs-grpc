package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"time"

	"scaling-up-rest-vs-grpc/internal/data/cache"
	"scaling-up-rest-vs-grpc/internal/grpc"
	"scaling-up-rest-vs-grpc/internal/lib/redis"
)

const listenAddr = ":50051"

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

	server := grpc.NewServer(cached)

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	slog.Info("grpc-server started", "addr", listenAddr)
	if err := server.Serve(lis); err != nil {
		slog.Error("grpc server stopped", "error", err)
		os.Exit(1)
	}
}
