package handler

import (
	"io"
	"log/slog"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/rest/service"
)

// Handler serves student data over REST in either JSON or Protobuf binary format, and contains no business logic beyond asking the service for data and writing the response.
type Handler struct {
	service *service.Service
}

// New creates a Handler backed by the given Service.
func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

// ServeJSONSmall writes the 1-entry dataset as JSON.
func (h *Handler) ServeJSONSmall(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetSmallDataset())
}

// ServeJSONMedium writes the 100-entry dataset as JSON.
func (h *Handler) ServeJSONMedium(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetMediumDataset())
}

// ServeJSONLarge writes the 1000-entry dataset as JSON.
func (h *Handler) ServeJSONLarge(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.service.GetLargeDataset())
}

// ServeProtobufSmall writes the 1-entry dataset as Protobuf binary.
func (h *Handler) ServeProtobufSmall(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetSmallDataset())
}

// ServeProtobufMedium writes the 100-entry dataset as Protobuf binary.
func (h *Handler) ServeProtobufMedium(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetMediumDataset())
}

// ServeProtobufLarge writes the 1000-entry dataset as Protobuf binary.
func (h *Handler) ServeProtobufLarge(w http.ResponseWriter, r *http.Request) {
	writeProtobuf(w, h.service.GetLargeDataset())
}

// CreateStudentJSON reads a single Student record from the request body, encoded as JSON, and appends it to the cache, the reverse of what the ServeJSON* handlers above measure.
func (h *Handler) CreateStudentJSON(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var student model.Student
	if err := protojson.Unmarshal(body, &student); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	h.service.CreateStudent(&student)
	writeCreated(w)
}

// CreateStudentProtobuf reads a single Student record from the request body, encoded as Protobuf binary, and appends it to the cache, the reverse of what the ServeProtobuf* handlers above measure.
func (h *Handler) CreateStudentProtobuf(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var student model.Student
	if err := proto.Unmarshal(body, &student); err != nil {
		http.Error(w, "invalid protobuf body", http.StatusBadRequest)
		return
	}

	h.service.CreateStudent(&student)
	writeCreated(w)
}

// writeJSON marshals data as JSON and writes it to w, or responds with an error if data is nil or marshaling fails.
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

// writeProtobuf marshals data as Protobuf binary and writes it to w, or responds with an error if data is nil or marshaling fails.
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

// writeCreated writes the minimal 201 Created acknowledgement shared by both CreateStudentJSON and CreateStudentProtobuf.
func writeCreated(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"created"}`))
}
