package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	"scaling-up-rest-vs-grpc/internal/data/constant"
	"scaling-up-rest-vs-grpc/internal/lib/redis"
	"scaling-up-rest-vs-grpc/internal/seeder"
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

	if err := seedDataset(client, constant.RedisKeySmallDataset, 1); err != nil {
		slog.Error("failed to seed small dataset", "error", err)
		os.Exit(1)
	}
	if err := seedDataset(client, constant.RedisKeyMediumDataset, 100); err != nil {
		slog.Error("failed to seed medium dataset", "error", err)
		os.Exit(1)
	}
	if err := seedDataset(client, constant.RedisKeyLargeDataset, 1000); err != nil {
		slog.Error("failed to seed large dataset", "error", err)
		os.Exit(1)
	}

	slog.Info("seeding completed", "redis_addr", addr)
}

// seedDataset generates n fake students, encodes them as a single StudentResponse, and writes the result to Redis under key.
func seedDataset(client *redis.Client, key string, n int) error {
	resp, err := seeder.ToStudentResponse(n)
	if err != nil {
		return err
	}
	b, err := protojson.Marshal(resp)
	if err != nil {
		return err
	}
	return client.Set(context.Background(), key, string(b))
}
