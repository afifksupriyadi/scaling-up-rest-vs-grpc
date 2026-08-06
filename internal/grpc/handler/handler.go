package handler

import (
	"context"

	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/grpc/service"
)

// Handler implements the generated StudentServiceServer interface, and relies on gRPC's own machinery for Protobuf serialization, so it only asks the service for data and returns it.
type Handler struct {
	model.UnimplementedStudentServiceServer
	service *service.Service
}

// New creates a Handler backed by the given Service.
func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

// GetStudentsSmall returns the 1-entry dataset.
func (h *Handler) GetStudentsSmall(ctx context.Context, _ *model.Empty) (*model.StudentResponse, error) {
	return h.service.GetSmallDataset(), nil
}

// GetStudentsMedium returns the 100-entry dataset.
func (h *Handler) GetStudentsMedium(ctx context.Context, _ *model.Empty) (*model.StudentResponse, error) {
	return h.service.GetMediumDataset(), nil
}

// GetStudentsLarge returns the 1000-entry dataset.
func (h *Handler) GetStudentsLarge(ctx context.Context, _ *model.Empty) (*model.StudentResponse, error) {
	return h.service.GetLargeDataset(), nil
}

// CreateStudent receives a single Student message as the request body itself, the reverse of the three GetStudents* methods above, and measures server-side deserialization cost of a Protobuf-encoded write.
func (h *Handler) CreateStudent(ctx context.Context, student *model.Student) (*model.CreateStudentResponse, error) {
	h.service.CreateStudent(student)
	return &model.CreateStudentResponse{Success: true}, nil
}
