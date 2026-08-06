package service

import (
	"scaling-up-rest-vs-grpc/internal/data/cache"
	"scaling-up-rest-vs-grpc/internal/data/model"
)

// Service exposes read-only access to the seeded student dataset for REST handlers, and contains no logic beyond delegating to Cache.
type Service struct {
	cache *cache.Cache
}

// New creates a Service backed by the given Cache.
func New(c *cache.Cache) *Service {
	return &Service{cache: c}
}

// GetSmallDataset returns the 1-entry dataset from the cache.
func (s *Service) GetSmallDataset() *model.StudentResponse {
	return s.cache.GetSmallDataset()
}

// GetMediumDataset returns the 100-entry dataset from the cache.
func (s *Service) GetMediumDataset() *model.StudentResponse {
	return s.cache.GetMediumDataset()
}

// GetLargeDataset returns the 1000-entry dataset from the cache.
func (s *Service) GetLargeDataset() *model.StudentResponse {
	return s.cache.GetLargeDataset()
}

// CreateStudent appends a new student record to the cache, the write-path counterpart of the three GetXDataset methods above.
func (s *Service) CreateStudent(student *model.Student) {
	s.cache.AddStudent(student)
}
