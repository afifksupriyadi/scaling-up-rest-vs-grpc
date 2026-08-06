package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"scaling-up-rest-vs-grpc/internal/data/cache"
	"scaling-up-rest-vs-grpc/internal/grpc"
	"scaling-up-rest-vs-grpc/internal/lib/redis"
)

const (
	listenAddr = ":50051"
	pprofAddr  = ":6060"
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

	// pprof runs on a separate port from the gRPC server itself, so it can
	// be profiled from outside without touching the gRPC handler code at all.
	go func() {
		slog.Info("pprof endpoint started", "addr", pprofAddr)
		http.ListenAndServe(pprofAddr, nil)
	}()

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
