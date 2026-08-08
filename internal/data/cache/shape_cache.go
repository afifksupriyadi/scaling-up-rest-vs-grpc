package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"

	"scaling-up-rest-vs-grpc/internal/data/constant"
	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/lib/redis"
)

// ShapeCache keeps all eight shape-experiment datasets (four structural
// depths x two calibrated size tiers) in memory, separate from Cache,
// so this investigative experiment never shares state with the Student
// dataset used by the main thesis benchmark.
type ShapeCache struct {
	mu sync.RWMutex

	depth0Compact *model.ShapeDepth0Response
	depth0Large   *model.ShapeDepth0Response

	depth1WideCompact *model.ShapeDepth1WideResponse
	depth1WideLarge   *model.ShapeDepth1WideResponse

	depth3NarrowCompact *model.ShapeDepth3NarrowResponse
	depth3NarrowLarge   *model.ShapeDepth3NarrowResponse

	depth4WideCompact *model.ShapeDepth4WideResponse
	depth4WideLarge   *model.ShapeDepth4WideResponse
}

// NewShape creates an empty ShapeCache with none of its datasets loaded yet.
func NewShape() *ShapeCache {
	return &ShapeCache{}
}

// LoadFromRedis reads all eight shape-experiment dataset keys from Redis and copies them into memory, and should be called once at server startup before serving any request.
func (c *ShapeCache) LoadFromRedis(ctx context.Context, client *redis.Client) error {
	d0c, err := readShapeDepth0FromRedis(ctx, client, constant.RedisKeyShapeDepth0Compact)
	if err != nil {
		return fmt.Errorf("load shape depth0 compact: %w", err)
	}
	d0l, err := readShapeDepth0FromRedis(ctx, client, constant.RedisKeyShapeDepth0Large)
	if err != nil {
		return fmt.Errorf("load shape depth0 large: %w", err)
	}
	d1c, err := readShapeDepth1WideFromRedis(ctx, client, constant.RedisKeyShapeDepth1WideCompact)
	if err != nil {
		return fmt.Errorf("load shape depth1-wide compact: %w", err)
	}
	d1l, err := readShapeDepth1WideFromRedis(ctx, client, constant.RedisKeyShapeDepth1WideLarge)
	if err != nil {
		return fmt.Errorf("load shape depth1-wide large: %w", err)
	}
	d3c, err := readShapeDepth3NarrowFromRedis(ctx, client, constant.RedisKeyShapeDepth3NarrowCompact)
	if err != nil {
		return fmt.Errorf("load shape depth3-narrow compact: %w", err)
	}
	d3l, err := readShapeDepth3NarrowFromRedis(ctx, client, constant.RedisKeyShapeDepth3NarrowLarge)
	if err != nil {
		return fmt.Errorf("load shape depth3-narrow large: %w", err)
	}
	d4c, err := readShapeDepth4WideFromRedis(ctx, client, constant.RedisKeyShapeDepth4WideCompact)
	if err != nil {
		return fmt.Errorf("load shape depth4-wide compact: %w", err)
	}
	d4l, err := readShapeDepth4WideFromRedis(ctx, client, constant.RedisKeyShapeDepth4WideLarge)
	if err != nil {
		return fmt.Errorf("load shape depth4-wide large: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.depth0Compact, c.depth0Large = d0c, d0l
	c.depth1WideCompact, c.depth1WideLarge = d1c, d1l
	c.depth3NarrowCompact, c.depth3NarrowLarge = d3c, d3l
	c.depth4WideCompact, c.depth4WideLarge = d4c, d4l
	return nil
}

func readShapeDepth0FromRedis(ctx context.Context, client *redis.Client, key string) (*model.ShapeDepth0Response, error) {
	raw, err := client.Get(ctx, key)
	if errors.Is(err, redis.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var resp model.ShapeDepth0Response
	if err := protojson.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal dataset: %w", err)
	}
	return &resp, nil
}

func readShapeDepth1WideFromRedis(ctx context.Context, client *redis.Client, key string) (*model.ShapeDepth1WideResponse, error) {
	raw, err := client.Get(ctx, key)
	if errors.Is(err, redis.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var resp model.ShapeDepth1WideResponse
	if err := protojson.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal dataset: %w", err)
	}
	return &resp, nil
}

func readShapeDepth3NarrowFromRedis(ctx context.Context, client *redis.Client, key string) (*model.ShapeDepth3NarrowResponse, error) {
	raw, err := client.Get(ctx, key)
	if errors.Is(err, redis.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var resp model.ShapeDepth3NarrowResponse
	if err := protojson.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal dataset: %w", err)
	}
	return &resp, nil
}

func readShapeDepth4WideFromRedis(ctx context.Context, client *redis.Client, key string) (*model.ShapeDepth4WideResponse, error) {
	raw, err := client.Get(ctx, key)
	if errors.Is(err, redis.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var resp model.ShapeDepth4WideResponse
	if err := protojson.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal dataset: %w", err)
	}
	return &resp, nil
}

// GetDepth0Compact returns the ~100 KB flat dataset already sitting in memory.
func (c *ShapeCache) GetDepth0Compact() *model.ShapeDepth0Response {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depth0Compact
}

// GetDepth0Large returns the ~500 KB flat dataset already sitting in memory.
func (c *ShapeCache) GetDepth0Large() *model.ShapeDepth0Response {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depth0Large
}

// GetDepth1WideCompact returns the ~100 KB depth-1-wide dataset already sitting in memory.
func (c *ShapeCache) GetDepth1WideCompact() *model.ShapeDepth1WideResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depth1WideCompact
}

// GetDepth1WideLarge returns the ~500 KB depth-1-wide dataset already sitting in memory.
func (c *ShapeCache) GetDepth1WideLarge() *model.ShapeDepth1WideResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depth1WideLarge
}

// GetDepth3NarrowCompact returns the ~100 KB depth-3-narrow dataset already sitting in memory.
func (c *ShapeCache) GetDepth3NarrowCompact() *model.ShapeDepth3NarrowResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depth3NarrowCompact
}

// GetDepth3NarrowLarge returns the ~500 KB depth-3-narrow dataset already sitting in memory.
func (c *ShapeCache) GetDepth3NarrowLarge() *model.ShapeDepth3NarrowResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depth3NarrowLarge
}

// GetDepth4WideCompact returns the ~100 KB depth-4-wide dataset already sitting in memory.
func (c *ShapeCache) GetDepth4WideCompact() *model.ShapeDepth4WideResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depth4WideCompact
}

// GetDepth4WideLarge returns the ~500 KB depth-4-wide dataset already sitting in memory.
func (c *ShapeCache) GetDepth4WideLarge() *model.ShapeDepth4WideResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depth4WideLarge
}
