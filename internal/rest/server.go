package rest

import (
	"net/http"

	"scaling-up-rest-vs-grpc/internal/data/store"
	"scaling-up-rest-vs-grpc/internal/rest/handler"
	"scaling-up-rest-vs-grpc/internal/rest/service"
)

// NewServer builds the REST http.Server with all 8 routes registered.
// The http1/http2 segments in the path are labels for JMeter Thread Group
// organization only; both route to the same handler, since HTTP/1.1 vs
// HTTP/2 is decided by ALPN negotiation, not by the URL.
func NewServer(addr string, dataStore *store.Store) *http.Server {
	svc := service.New(dataStore)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /rest/json/http1/students/small", h.ServeJSONSmall)
	mux.HandleFunc("GET /rest/json/http1/students/large", h.ServeJSONLarge)
	mux.HandleFunc("GET /rest/json/http2/students/small", h.ServeJSONSmall)
	mux.HandleFunc("GET /rest/json/http2/students/large", h.ServeJSONLarge)
	mux.HandleFunc("GET /rest/protobuf/http1/students/small", h.ServeProtobufSmall)
	mux.HandleFunc("GET /rest/protobuf/http1/students/large", h.ServeProtobufLarge)
	mux.HandleFunc("GET /rest/protobuf/http2/students/small", h.ServeProtobufSmall)
	mux.HandleFunc("GET /rest/protobuf/http2/students/large", h.ServeProtobufLarge)

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
