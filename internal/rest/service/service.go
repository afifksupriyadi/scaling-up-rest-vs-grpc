package service

import (
	"scaling-up-rest-vs-grpc/internal/data/cache"
	"scaling-up-rest-vs-grpc/internal/data/model"
)

// Service exposes read-only access to the seeded order dataset for REST handlers, and contains no logic beyond delegating to Cache.
type Service struct {
	cache *cache.Cache
}

// New creates a Service backed by the given Cache.
func New(c *cache.Cache) *Service {
	return &Service{cache: c}
}

// GetDepthZero returns the depth-0 dataset from the cache.
func (s *Service) GetDepthZero() *model.OrderDepthZeroResponse {
	return s.cache.GetDepthZero()
}

// GetDepthTwo returns the depth-2 dataset from the cache.
func (s *Service) GetDepthTwo() *model.OrderDepthTwoResponse {
	return s.cache.GetDepthTwo()
}

// GetDepthFour returns the depth-4 dataset from the cache.
func (s *Service) GetDepthFour() *model.OrderDepthFourResponse {
	return s.cache.GetDepthFour()
}

// GetOne returns the 1-element dataset from the cache.
func (s *Service) GetOne() *model.OrderDepthZeroResponse {
	return s.cache.GetOne()
}

// GetHundred returns the 100-element dataset from the cache.
func (s *Service) GetHundred() *model.OrderDepthZeroResponse {
	return s.cache.GetHundred()
}

// GetThousand returns the 1000-element dataset from the cache.
func (s *Service) GetThousand() *model.OrderDepthZeroResponse {
	return s.cache.GetThousand()
}
