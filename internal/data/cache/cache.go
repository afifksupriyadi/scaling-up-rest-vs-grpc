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

// Cache keeps the seeded student data in memory across all three sizes so request handlers never touch Redis while the benchmark is running.
type Cache struct {
	mu     sync.RWMutex
	small  *model.StudentResponse
	medium *model.StudentResponse
	large  *model.StudentResponse
	posted []*model.Student
}

// New creates an empty Cache with none of its datasets loaded yet.
func New() *Cache {
	return &Cache{}
}

// LoadFromRedis reads all three dataset keys from Redis and copies them into memory, and should be called once at server startup before serving any request.
func (c *Cache) LoadFromRedis(ctx context.Context, client *redis.Client) error {
	small, err := readDatasetFromRedis(ctx, client, constant.RedisKeySmallDataset)
	if err != nil {
		return fmt.Errorf("load small dataset: %w", err)
	}
	medium, err := readDatasetFromRedis(ctx, client, constant.RedisKeyMediumDataset)
	if err != nil {
		return fmt.Errorf("load medium dataset: %w", err)
	}
	large, err := readDatasetFromRedis(ctx, client, constant.RedisKeyLargeDataset)
	if err != nil {
		return fmt.Errorf("load large dataset: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.small = small
	c.medium = medium
	c.large = large
	return nil
}

// readDatasetFromRedis reads and decodes a single dataset key from Redis, returning a nil result without error if the key does not exist yet.
func readDatasetFromRedis(ctx context.Context, client *redis.Client, key string) (*model.StudentResponse, error) {
	raw, err := client.Get(ctx, key)
	if errors.Is(err, redis.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var resp model.StudentResponse
	if err := protojson.Unmarshal([]byte(raw), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal dataset: %w", err)
	}
	return &resp, nil
}

// GetSmallDataset returns the 1-entry dataset already sitting in memory, does not touch Redis, and returns nil if LoadFromRedis has not been called yet.
func (c *Cache) GetSmallDataset() *model.StudentResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.small
}

// GetMediumDataset returns the 100-entry dataset already sitting in memory, does not touch Redis, and returns nil if LoadFromRedis has not been called yet.
func (c *Cache) GetMediumDataset() *model.StudentResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.medium
}

// GetLargeDataset returns the 1000-entry dataset already sitting in memory, does not touch Redis, and returns nil if LoadFromRedis has not been called yet.
func (c *Cache) GetLargeDataset() *model.StudentResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.large
}

// AddStudent appends a student record submitted via POST, kept in a slice separate from small/medium/large so write operations never alter the fixed-size datasets used by the GET benchmarks.
// - Uses the write side of the mutex, since this mutates state that concurrent readers may access.
func (c *Cache) AddStudent(s *model.Student) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.posted = append(c.posted, s)
}
