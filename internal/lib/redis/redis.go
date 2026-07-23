package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrNotFound is returned by Get when the key does not exist in Redis.
var ErrNotFound = errors.New("key not found")

// Client wraps go-redis so the rest of the codebase depends on this
// package instead of importing go-redis directly everywhere.
type Client struct {
	rdb *redis.Client
}

// NewClient creates a Redis client connected to the given address (host:port).
func NewClient(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Client{rdb: rdb}
}

// Ping verifies the connection to Redis is alive.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// WaitReady retries Ping until it succeeds or the attempt budget runs out.
// Needed because Redis may not be reachable yet right after container startup.
func (c *Client) WaitReady(ctx context.Context, attempts int, delay time.Duration) error {
	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := c.Ping(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("redis not ready after %d attempts: %w", attempts, lastErr)
}

// Set stores a string value under the given key.
func (c *Client) Set(ctx context.Context, key string, value string) error {
	return c.rdb.Set(ctx, key, value, 0).Err()
}

// Get retrieves the string value stored under the given key. It returns
// ErrNotFound if the key does not exist.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get key %s: %w", key, err)
	}
	return val, nil
}
