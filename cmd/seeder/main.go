package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

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

	if err := seedOrderExperiment(client); err != nil {
		slog.Error("failed to seed order experiment datasets", "error", err)
		os.Exit(1)
	}

	if err := seedShapeExperiment(client); err != nil {
		slog.Error("failed to seed shape experiment datasets", "error", err)
		os.Exit(1)
	}

	slog.Info("seeding completed", "redis_addr", addr)
}

// seedOrderExperiment seeds all six datasets used by the depth-scenario
// and element-count-scenario investigation, in two groups: depth
// (DepthZero, DepthTwo, DepthFour) followed by element count
// (One, Hundred, Thousand). Every dataset follows the same
// generate-then-write pattern, regardless of which group it belongs to.
func seedOrderExperiment(client *redis.Client) error {
	// Depth scenario.
	depthZero, err := seeder.ToOrderDepthZeroResponse(1)
	if err != nil {
		return err
	}
	if err := writeOrderDataset(client, constant.RedisKeyOrderDepthZero, depthZero); err != nil {
		return err
	}

	depthTwo, err := seeder.ToOrderDepthTwoResponse(1)
	if err != nil {
		return err
	}
	if err := writeOrderDataset(client, constant.RedisKeyOrderDepthTwo, depthTwo); err != nil {
		return err
	}

	depthFour, err := seeder.ToOrderDepthFourResponse(1)
	if err != nil {
		return err
	}
	if err := writeOrderDataset(client, constant.RedisKeyOrderDepthFour, depthFour); err != nil {
		return err
	}

	// Element-count scenario.
	one, err := seeder.ToOrderDepthZeroResponse(1)
	if err != nil {
		return err
	}
	if err := writeOrderDataset(client, constant.RedisKeyOrderOne, one); err != nil {
		return err
	}

	hundred, err := seeder.ToOrderDepthZeroResponse(100)
	if err != nil {
		return err
	}
	if err := writeOrderDataset(client, constant.RedisKeyOrderHundred, hundred); err != nil {
		return err
	}

	thousand, err := seeder.ToOrderDepthZeroResponse(1000)
	if err != nil {
		return err
	}
	return writeOrderDataset(client, constant.RedisKeyOrderThousand, thousand)
}

// writeOrderDataset encodes a single order-experiment message as protojson and writes it to Redis under key.
func writeOrderDataset(client *redis.Client, key string, msg proto.Message) error {
	b, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}
	return client.Set(context.Background(), key, string(b))
}

// seedShapeExperiment seeds all eight datasets (four structural-depth
// variants x two calibrated size tiers) used by the shape-experiment
// investigation. Document counts per key were derived from cmd/calibrate-shape,
// which measures actual Protobuf-marshaled size rather than estimating it.
func seedShapeExperiment(client *redis.Client) error {
	if err := seedShapeDataset(client, constant.RedisKeyShapeDepth0Compact, seeder.ToShapeDepth0Response(87)); err != nil {
		return err
	}
	if err := seedShapeDataset(client, constant.RedisKeyShapeDepth0Large, seeder.ToShapeDepth0Response(437)); err != nil {
		return err
	}
	if err := seedShapeDataset(client, constant.RedisKeyShapeDepth1WideCompact, seeder.ToShapeDepth1WideResponse(33)); err != nil {
		return err
	}
	if err := seedShapeDataset(client, constant.RedisKeyShapeDepth1WideLarge, seeder.ToShapeDepth1WideResponse(163)); err != nil {
		return err
	}
	if err := seedShapeDataset(client, constant.RedisKeyShapeDepth3NarrowCompact, seeder.ToShapeDepth3NarrowResponse(11)); err != nil {
		return err
	}
	if err := seedShapeDataset(client, constant.RedisKeyShapeDepth3NarrowLarge, seeder.ToShapeDepth3NarrowResponse(57)); err != nil {
		return err
	}
	if err := seedShapeDataset(client, constant.RedisKeyShapeDepth4WideCompact, seeder.ToShapeDepth4WideResponse(6)); err != nil {
		return err
	}
	if err := seedShapeDataset(client, constant.RedisKeyShapeDepth4WideLarge, seeder.ToShapeDepth4WideResponse(28)); err != nil {
		return err
	}
	return nil
}

// seedShapeDataset encodes a single shape-experiment message as protojson and writes it to Redis under key.
func seedShapeDataset(client *redis.Client, key string, msg proto.Message) error {
	b, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}
	return client.Set(context.Background(), key, string(b))
}
