package handler

import (
	"log/slog"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/rest/service"
)

// Handler serves student data over REST in either JSON or Protobuf binary
// format. It contains no business logic: it only asks the service for
// data and writes the response.
type Handler struct {
	service *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) ServeJSONSmall(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetSmallDataset())
}

func (h *Handler) ServeJSONLarge(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetLargeDataset())
}

func (h *Handler) ServeProtobufSmall(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetSmallDataset())
}

func (h *Handler) ServeProtobufLarge(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetLargeDataset())
}

func writeJSON(w http.ResponseWriter, data *model.StudentResponse) {
	if data == nil {
		http.Error(w, "dataset not seeded yet", http.StatusServiceUnavailable)
		return
	}
	b, err := protojson.Marshal(data)
	if err != nil {
		slog.Error("failed to marshal response as json", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func writeProtobuf(w http.ResponseWriter, data *model.StudentResponse) {
	if data == nil {
		http.Error(w, "dataset not seeded yet", http.StatusServiceUnavailable)
		return
	}
	b, err := proto.Marshal(data)
	if err != nil {
		slog.Error("failed to marshal response as protobuf", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Write(b)
}
