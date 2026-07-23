package handler

import (
	"context"

	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/grpc/service"
)

// Handler implements the generated StudentServiceServer interface. gRPC's
// own machinery handles Protobuf serialization, so this only asks the
// service for data and returns it.
type Handler struct {
	model.UnimplementedStudentServiceServer
	service *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetStudent(ctx context.Context, _ *model.Empty) (*model.StudentResponse, error) {
	return h.service.GetSmallDataset(), nil
}

func (h *Handler) GetStudents(ctx context.Context, _ *model.Empty) (*model.StudentResponse, error) {
	return h.service.GetLargeDataset(), nil
}
