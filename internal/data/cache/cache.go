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

// Cache keeps the seeded order data in memory across all six datasets
// (three depth-scenario points, three element-count-scenario points) so
// request handlers never touch Redis while the benchmark is running.
type Cache struct {
	mu sync.RWMutex

	// Depth scenario.
	depthZero *model.OrderDepthZeroResponse
	depthTwo  *model.OrderDepthTwoResponse
	depthFour *model.OrderDepthFourResponse

	// Element-count scenario.
	one      *model.OrderDepthZeroResponse
	hundred  *model.OrderDepthZeroResponse
	thousand *model.OrderDepthZeroResponse
}

// New creates an empty Cache with none of its datasets loaded yet.
func New() *Cache {
	return &Cache{}
}

// LoadFromRedis reads all six dataset keys from Redis and copies them into memory, and should be called once at server startup before serving any request.
func (c *Cache) LoadFromRedis(ctx context.Context, client *redis.Client) error {
	depthZero, err := readOrderDepthZeroFromRedis(ctx, client, constant.RedisKeyOrderDepthZero)
	if err != nil {
		return fmt.Errorf("load order depth-zero: %w", err)
	}
	depthTwo, err := readOrderDepthTwoFromRedis(ctx, client, constant.RedisKeyOrderDepthTwo)
	if err != nil {
		return fmt.Errorf("load order depth-two: %w", err)
	}
	depthFour, err := readOrderDepthFourFromRedis(ctx, client, constant.RedisKeyOrderDepthFour)
	if err != nil {
		return fmt.Errorf("load order depth-four: %w", err)
	}
	one, err := readOrderDepthZeroFromRedis(ctx, client, constant.RedisKeyOrderOne)
	if err != nil {
		return fmt.Errorf("load order one: %w", err)
	}
	hundred, err := readOrderDepthZeroFromRedis(ctx, client, constant.RedisKeyOrderHundred)
	if err != nil {
		return fmt.Errorf("load order hundred: %w", err)
	}
	thousand, err := readOrderDepthZeroFromRedis(ctx, client, constant.RedisKeyOrderThousand)
	if err != nil {
		return fmt.Errorf("load order thousand: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.depthZero, c.depthTwo, c.depthFour = depthZero, depthTwo, depthFour
	c.one, c.hundred, c.thousand = one, hundred, thousand
	return nil
}

func readOrderDepthZeroFromRedis(ctx context.Context, client *redis.Client, key string) (*model.OrderDepthZeroResponse, error) {
	raw, err := client.Get(ctx, key)
	if errors.Is(err, redis.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var resp model.OrderDepthZeroResponse
	if err := protojson.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal dataset: %w", err)
	}
	return &resp, nil
}

func readOrderDepthTwoFromRedis(ctx context.Context, client *redis.Client, key string) (*model.OrderDepthTwoResponse, error) {
	raw, err := client.Get(ctx, key)
	if errors.Is(err, redis.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var resp model.OrderDepthTwoResponse
	if err := protojson.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal dataset: %w", err)
	}
	return &resp, nil
}

func readOrderDepthFourFromRedis(ctx context.Context, client *redis.Client, key string) (*model.OrderDepthFourResponse, error) {
	raw, err := client.Get(ctx, key)
	if errors.Is(err, redis.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var resp model.OrderDepthFourResponse
	if err := protojson.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal dataset: %w", err)
	}
	return &resp, nil
}

// GetDepthZero returns the depth-0 dataset already sitting in memory.
func (c *Cache) GetDepthZero() *model.OrderDepthZeroResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depthZero
}

// GetDepthTwo returns the depth-2 dataset already sitting in memory.
func (c *Cache) GetDepthTwo() *model.OrderDepthTwoResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depthTwo
}

// GetDepthFour returns the depth-4 dataset already sitting in memory.
func (c *Cache) GetDepthFour() *model.OrderDepthFourResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.depthFour
}

// GetOne returns the 1-element dataset already sitting in memory.
func (c *Cache) GetOne() *model.OrderDepthZeroResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.one
}

// GetHundred returns the 100-element dataset already sitting in memory.
func (c *Cache) GetHundred() *model.OrderDepthZeroResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hundred
}

// GetThousand returns the 1000-element dataset already sitting in memory.
func (c *Cache) GetThousand() *model.OrderDepthZeroResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.thousand
}
