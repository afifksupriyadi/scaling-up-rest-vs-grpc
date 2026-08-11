package handler

import (
	"context"

	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/grpc/service"
)

// Handler implements the generated OrderExperimentServiceServer interface,
// and relies on gRPC's own machinery for Protobuf serialization, so it
// only asks the service for data and returns it.
type Handler struct {
	model.UnimplementedOrderExperimentServiceServer
	service *service.Service
}

// New creates a Handler backed by the given Service.
func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

// GetOrderDepthZero returns the depth-0 dataset.
func (h *Handler) GetOrderDepthZero(ctx context.Context, _ *model.OrderEmpty) (*model.OrderDepthZeroResponse, error) {
	return h.service.GetDepthZero(), nil
}

// GetOrderDepthTwo returns the depth-2 dataset.
func (h *Handler) GetOrderDepthTwo(ctx context.Context, _ *model.OrderEmpty) (*model.OrderDepthTwoResponse, error) {
	return h.service.GetDepthTwo(), nil
}

// GetOrderDepthFour returns the depth-4 dataset.
func (h *Handler) GetOrderDepthFour(ctx context.Context, _ *model.OrderEmpty) (*model.OrderDepthFourResponse, error) {
	return h.service.GetDepthFour(), nil
}

// GetOrderOne returns the 1-element dataset.
func (h *Handler) GetOrderOne(ctx context.Context, _ *model.OrderEmpty) (*model.OrderDepthZeroResponse, error) {
	return h.service.GetOne(), nil
}

// GetOrderHundred returns the 100-element dataset.
func (h *Handler) GetOrderHundred(ctx context.Context, _ *model.OrderEmpty) (*model.OrderDepthZeroResponse, error) {
	return h.service.GetHundred(), nil
}

// GetOrderThousand returns the 1000-element dataset.
func (h *Handler) GetOrderThousand(ctx context.Context, _ *model.OrderEmpty) (*model.OrderDepthZeroResponse, error) {
	return h.service.GetThousand(), nil
}
