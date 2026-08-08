package service

import (
	"scaling-up-rest-vs-grpc/internal/data/cache"
	"scaling-up-rest-vs-grpc/internal/data/model"
)

// ShapeService exposes read-only access to the eight shape-experiment
// datasets for REST handlers, and contains no logic beyond delegating
// to ShapeCache.
type ShapeService struct {
	cache *cache.ShapeCache
}

// NewShape creates a ShapeService backed by the given ShapeCache.
func NewShape(c *cache.ShapeCache) *ShapeService {
	return &ShapeService{cache: c}
}

// GetDepth0Compact returns the ~100 KB flat dataset from the cache.
func (s *ShapeService) GetDepth0Compact() *model.ShapeDepth0Response {
	return s.cache.GetDepth0Compact()
}

// GetDepth0Large returns the ~500 KB flat dataset from the cache.
func (s *ShapeService) GetDepth0Large() *model.ShapeDepth0Response {
	return s.cache.GetDepth0Large()
}

// GetDepth1WideCompact returns the ~100 KB depth-1-wide dataset from the cache.
func (s *ShapeService) GetDepth1WideCompact() *model.ShapeDepth1WideResponse {
	return s.cache.GetDepth1WideCompact()
}

// GetDepth1WideLarge returns the ~500 KB depth-1-wide dataset from the cache.
func (s *ShapeService) GetDepth1WideLarge() *model.ShapeDepth1WideResponse {
	return s.cache.GetDepth1WideLarge()
}

// GetDepth3NarrowCompact returns the ~100 KB depth-3-narrow dataset from the cache.
func (s *ShapeService) GetDepth3NarrowCompact() *model.ShapeDepth3NarrowResponse {
	return s.cache.GetDepth3NarrowCompact()
}

// GetDepth3NarrowLarge returns the ~500 KB depth-3-narrow dataset from the cache.
func (s *ShapeService) GetDepth3NarrowLarge() *model.ShapeDepth3NarrowResponse {
	return s.cache.GetDepth3NarrowLarge()
}

// GetDepth4WideCompact returns the ~100 KB depth-4-wide dataset from the cache.
func (s *ShapeService) GetDepth4WideCompact() *model.ShapeDepth4WideResponse {
	return s.cache.GetDepth4WideCompact()
}

// GetDepth4WideLarge returns the ~500 KB depth-4-wide dataset from the cache.
func (s *ShapeService) GetDepth4WideLarge() *model.ShapeDepth4WideResponse {
	return s.cache.GetDepth4WideLarge()
}
