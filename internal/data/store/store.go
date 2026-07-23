package store

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

// Store keeps the seeded student data in memory so request handlers never
// touch Redis while the benchmark is running.
type Store struct {
	mu    sync.RWMutex
	small *model.StudentResponse
	large *model.StudentResponse
}

func New() *Store {
	return &Store{}
}

// LoadFromRedis reads both dataset keys from Redis and copies them into
// memory. Call this once, at server startup, before serving any request.
func (s *Store) LoadFromRedis(ctx context.Context, client *redis.Client) error {
	small, err := readDatasetFromRedis(ctx, client, constant.RedisKeySmallDataset)
	if err != nil {
		return fmt.Errorf("load small dataset: %w", err)
	}
	large, err := readDatasetFromRedis(ctx, client, constant.RedisKeyLargeDataset)
	if err != nil {
		return fmt.Errorf("load large dataset: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.small = small
	s.large = large
	return nil
}

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

// GetSmallDataset returns the 1-entry dataset already sitting in memory.
// Does not touch Redis. Call LoadFromRedis first, or this returns nil.
func (s *Store) GetSmallDataset() *model.StudentResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.small
}

// GetLargeDataset returns the 100-entry dataset already sitting in memory.
// Does not touch Redis. Call LoadFromRedis first, or this returns nil.
func (s *Store) GetLargeDataset() *model.StudentResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.large
}
