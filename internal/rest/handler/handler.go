package handler

import (
	"log/slog"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"scaling-up-rest-vs-grpc/internal/rest/service"
)

// Handler serves order data over REST in either JSON or Protobuf binary
// format, and contains no business logic beyond asking the service for
// data and writing the response.
type Handler struct {
	service *service.Service
}

// New creates a Handler backed by the given Service.
func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

// ServeJSONDepthZero writes the depth-0 dataset as JSON.
func (h *Handler) ServeJSONDepthZero(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetDepthZero())
}

// ServeJSONDepthTwo writes the depth-2 dataset as JSON.
func (h *Handler) ServeJSONDepthTwo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetDepthTwo())
}

// ServeJSONDepthFour writes the depth-4 dataset as JSON.
func (h *Handler) ServeJSONDepthFour(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetDepthFour())
}

// ServeJSONOne writes the 1-element dataset as JSON.
func (h *Handler) ServeJSONOne(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetOne())
}

// ServeJSONHundred writes the 100-element dataset as JSON.
func (h *Handler) ServeJSONHundred(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetHundred())
}

// ServeJSONThousand writes the 1000-element dataset as JSON.
func (h *Handler) ServeJSONThousand(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetThousand())
}

// ServeProtobufDepthZero writes the depth-0 dataset as Protobuf binary.
func (h *Handler) ServeProtobufDepthZero(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetDepthZero())
}

// ServeProtobufDepthTwo writes the depth-2 dataset as Protobuf binary.
func (h *Handler) ServeProtobufDepthTwo(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetDepthTwo())
}

// ServeProtobufDepthFour writes the depth-4 dataset as Protobuf binary.
func (h *Handler) ServeProtobufDepthFour(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetDepthFour())
}

// ServeProtobufOne writes the 1-element dataset as Protobuf binary.
func (h *Handler) ServeProtobufOne(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetOne())
}

// ServeProtobufHundred writes the 100-element dataset as Protobuf binary.
func (h *Handler) ServeProtobufHundred(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetHundred())
}

// ServeProtobufThousand writes the 1000-element dataset as Protobuf binary.
func (h *Handler) ServeProtobufThousand(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetThousand())
}

// writeJSON marshals a response message as JSON and writes it to w.
// - Takes proto.Message (an interface) instead of a concrete type, since the three depth points are three distinct generated Go types.
func writeJSON(w http.ResponseWriter, data proto.Message) {
	b, err := protojson.Marshal(data)
	if err != nil {
		slog.Error("failed to marshal response as json", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// writeProtobuf marshals a response message as Protobuf binary and writes it to w.
// - Takes proto.Message (an interface) instead of a concrete type, since the three depth points are three distinct generated Go types.
func writeProtobuf(w http.ResponseWriter, data proto.Message) {
	b, err := proto.Marshal(data)
	if err != nil {
		slog.Error("failed to marshal response as protobuf", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Write(b)
}
