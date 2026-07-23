package service

import (
	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/data/store"
)

// Service exposes read-only access to the seeded student dataset for
// REST handlers. It contains no logic beyond delegating to Store.
type Service struct {
	store *store.Store
}

func New(s *store.Store) *Service {
	return &Service{store: s}
}

func (s *Service) GetSmallDataset() *model.StudentResponse {
	return s.store.GetSmallDataset()
}

func (s *Service) GetLargeDataset() *model.StudentResponse {
	return s.store.GetLargeDataset()
}
