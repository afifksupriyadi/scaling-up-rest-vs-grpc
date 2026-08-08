package handler

import (
	"log/slog"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"scaling-up-rest-vs-grpc/internal/rest/service"
)

// ShapeHandler serves the structural-depth investigation data over REST,
// in either JSON or Protobuf binary format. GET only; POST is
// intentionally out of scope for this experiment.
type ShapeHandler struct {
	service *service.ShapeService
}

// NewShape creates a ShapeHandler backed by the given ShapeService.
func NewShape(s *service.ShapeService) *ShapeHandler {
	return &ShapeHandler{service: s}
}

// ServeJSONDepth0Compact writes the ~100 KB flat dataset as JSON.
func (h *ShapeHandler) ServeJSONDepth0Compact(w http.ResponseWriter, r *http.Request) {
	writeShapeJSON(w, h.service.GetDepth0Compact())
}

// ServeJSONDepth0Large writes the ~500 KB flat dataset as JSON.
func (h *ShapeHandler) ServeJSONDepth0Large(w http.ResponseWriter, r *http.Request) {
	writeShapeJSON(w, h.service.GetDepth0Large())
}

// ServeJSONDepth1WideCompact writes the ~100 KB depth-1-wide dataset as JSON.
func (h *ShapeHandler) ServeJSONDepth1WideCompact(w http.ResponseWriter, r *http.Request) {
	writeShapeJSON(w, h.service.GetDepth1WideCompact())
}

// ServeJSONDepth1WideLarge writes the ~500 KB depth-1-wide dataset as JSON.
func (h *ShapeHandler) ServeJSONDepth1WideLarge(w http.ResponseWriter, r *http.Request) {
	writeShapeJSON(w, h.service.GetDepth1WideLarge())
}

// ServeJSONDepth3NarrowCompact writes the ~100 KB depth-3-narrow dataset as JSON.
func (h *ShapeHandler) ServeJSONDepth3NarrowCompact(w http.ResponseWriter, r *http.Request) {
	writeShapeJSON(w, h.service.GetDepth3NarrowCompact())
}

// ServeJSONDepth3NarrowLarge writes the ~500 KB depth-3-narrow dataset as JSON.
func (h *ShapeHandler) ServeJSONDepth3NarrowLarge(w http.ResponseWriter, r *http.Request) {
	writeShapeJSON(w, h.service.GetDepth3NarrowLarge())
}

// ServeJSONDepth4WideCompact writes the ~100 KB depth-4-wide dataset as JSON.
func (h *ShapeHandler) ServeJSONDepth4WideCompact(w http.ResponseWriter, r *http.Request) {
	writeShapeJSON(w, h.service.GetDepth4WideCompact())
}

// ServeJSONDepth4WideLarge writes the ~500 KB depth-4-wide dataset as JSON.
func (h *ShapeHandler) ServeJSONDepth4WideLarge(w http.ResponseWriter, r *http.Request) {
	writeShapeJSON(w, h.service.GetDepth4WideLarge())
}

// ServeProtobufDepth0Compact writes the ~100 KB flat dataset as Protobuf binary.
func (h *ShapeHandler) ServeProtobufDepth0Compact(w http.ResponseWriter, r *http.Request) {
	writeShapeProtobuf(w, h.service.GetDepth0Compact())
}

// ServeProtobufDepth0Large writes the ~500 KB flat dataset as Protobuf binary.
func (h *ShapeHandler) ServeProtobufDepth0Large(w http.ResponseWriter, r *http.Request) {
	writeShapeProtobuf(w, h.service.GetDepth0Large())
}

// ServeProtobufDepth1WideCompact writes the ~100 KB depth-1-wide dataset as Protobuf binary.
func (h *ShapeHandler) ServeProtobufDepth1WideCompact(w http.ResponseWriter, r *http.Request) {
	writeShapeProtobuf(w, h.service.GetDepth1WideCompact())
}

// ServeProtobufDepth1WideLarge writes the ~500 KB depth-1-wide dataset as Protobuf binary.
func (h *ShapeHandler) ServeProtobufDepth1WideLarge(w http.ResponseWriter, r *http.Request) {
	writeShapeProtobuf(w, h.service.GetDepth1WideLarge())
}

// ServeProtobufDepth3NarrowCompact writes the ~100 KB depth-3-narrow dataset as Protobuf binary.
func (h *ShapeHandler) ServeProtobufDepth3NarrowCompact(w http.ResponseWriter, r *http.Request) {
	writeShapeProtobuf(w, h.service.GetDepth3NarrowCompact())
}

// ServeProtobufDepth3NarrowLarge writes the ~500 KB depth-3-narrow dataset as Protobuf binary.
func (h *ShapeHandler) ServeProtobufDepth3NarrowLarge(w http.ResponseWriter, r *http.Request) {
	writeShapeProtobuf(w, h.service.GetDepth3NarrowLarge())
}

// ServeProtobufDepth4WideCompact writes the ~100 KB depth-4-wide dataset as Protobuf binary.
func (h *ShapeHandler) ServeProtobufDepth4WideCompact(w http.ResponseWriter, r *http.Request) {
	writeShapeProtobuf(w, h.service.GetDepth4WideCompact())
}

// ServeProtobufDepth4WideLarge writes the ~500 KB depth-4-wide dataset as Protobuf binary.
func (h *ShapeHandler) ServeProtobufDepth4WideLarge(w http.ResponseWriter, r *http.Request) {
	writeShapeProtobuf(w, h.service.GetDepth4WideLarge())
}

// writeShapeJSON marshals a shape-experiment response message as JSON and writes it to w.
// - Takes proto.Message (an interface) instead of a concrete type, since the four shape variants are four distinct generated Go types.
func writeShapeJSON(w http.ResponseWriter, data proto.Message) {
	b, err := protojson.Marshal(data)
	if err != nil {
		slog.Error("failed to marshal shape response as json", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// writeShapeProtobuf marshals a shape-experiment response message as Protobuf binary and writes it to w.
// - Takes proto.Message (an interface) instead of a concrete type, since the four shape variants are four distinct generated Go types.
func writeShapeProtobuf(w http.ResponseWriter, data proto.Message) {
	b, err := proto.Marshal(data)
	if err != nil {
		slog.Error("failed to marshal shape response as protobuf", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Write(b)
}
