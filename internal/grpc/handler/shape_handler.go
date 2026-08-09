package handler

import (
	"context"

	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/grpc/service"
)

// ShapeHandler implements the generated ShapeExperimentServiceServer
// interface, and relies on gRPC's own machinery for Protobuf
// serialization, so it only asks the service for data and returns it.
type ShapeHandler struct {
	model.UnimplementedShapeExperimentServiceServer
	service *service.ShapeService
}

// NewShape creates a ShapeHandler backed by the given ShapeService.
func NewShape(s *service.ShapeService) *ShapeHandler {
	return &ShapeHandler{service: s}
}

// GetShapeDepth0Compact returns the ~100 KB flat dataset.
func (h *ShapeHandler) GetShapeDepth0Compact(ctx context.Context, _ *model.ShapeEmpty) (*model.ShapeDepth0Response, error) {
	return h.service.GetDepth0Compact(), nil
}

// GetShapeDepth0Large returns the ~500 KB flat dataset.
func (h *ShapeHandler) GetShapeDepth0Large(ctx context.Context, _ *model.ShapeEmpty) (*model.ShapeDepth0Response, error) {
	return h.service.GetDepth0Large(), nil
}

// GetShapeDepth1WideCompact returns the ~100 KB depth-1-wide dataset.
func (h *ShapeHandler) GetShapeDepth1WideCompact(ctx context.Context, _ *model.ShapeEmpty) (*model.ShapeDepth1WideResponse, error) {
	return h.service.GetDepth1WideCompact(), nil
}

// GetShapeDepth1WideLarge returns the ~500 KB depth-1-wide dataset.
func (h *ShapeHandler) GetShapeDepth1WideLarge(ctx context.Context, _ *model.ShapeEmpty) (*model.ShapeDepth1WideResponse, error) {
	return h.service.GetDepth1WideLarge(), nil
}

// GetShapeDepth3NarrowCompact returns the ~100 KB depth-3-narrow dataset.
func (h *ShapeHandler) GetShapeDepth3NarrowCompact(ctx context.Context, _ *model.ShapeEmpty) (*model.ShapeDepth3NarrowResponse, error) {
	return h.service.GetDepth3NarrowCompact(), nil
}

// GetShapeDepth3NarrowLarge returns the ~500 KB depth-3-narrow dataset.
func (h *ShapeHandler) GetShapeDepth3NarrowLarge(ctx context.Context, _ *model.ShapeEmpty) (*model.ShapeDepth3NarrowResponse, error) {
	return h.service.GetDepth3NarrowLarge(), nil
}

// GetShapeDepth4WideCompact returns the ~100 KB depth-4-wide dataset.
func (h *ShapeHandler) GetShapeDepth4WideCompact(ctx context.Context, _ *model.ShapeEmpty) (*model.ShapeDepth4WideResponse, error) {
	return h.service.GetDepth4WideCompact(), nil
}

// GetShapeDepth4WideLarge returns the ~500 KB depth-4-wide dataset.
func (h *ShapeHandler) GetShapeDepth4WideLarge(ctx context.Context, _ *model.ShapeEmpty) (*model.ShapeDepth4WideResponse, error) {
	return h.service.GetDepth4WideLarge(), nil
}

// Ping returns an empty response immediately, used as a near-zero-cost reference point to isolate fixed per-call overhead from serialization cost.
func (h *ShapeHandler) Ping(ctx context.Context, _ *model.ShapeEmpty) (*model.ShapeEmpty, error) {
	return &model.ShapeEmpty{}, nil
}
