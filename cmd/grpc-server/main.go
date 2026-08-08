package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
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
	// Enables the block profiler, sampling every blocking event (channel,
	// mutex, etc.), so /debug/pprof/block returns meaningful data instead
	// of an empty profile, which is Go's default when this is unset.
	runtime.SetBlockProfileRate(1)

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

	shapeCached := cache.NewShape()
	if err := shapeCached.LoadFromRedis(context.Background(), client); err != nil {
		slog.Error("failed to load shape experiment dataset from redis", "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("pprof endpoint started", "addr", pprofAddr)
		http.ListenAndServe(pprofAddr, nil)
	}()

	server := grpc.NewServer(cached, shapeCached)

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
