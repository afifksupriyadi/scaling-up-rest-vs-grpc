package rest

import (
	"net/http"

	"scaling-up-rest-vs-grpc/internal/data/cache"
	"scaling-up-rest-vs-grpc/internal/rest/handler"
	"scaling-up-rest-vs-grpc/internal/rest/service"
)

// NewServers builds two REST http.Server instances sharing the same handler, one restricted to plain HTTP/1.1 on http1Addr and one restricted to unencrypted HTTP/2 (h2c) on http2Addr.
// - Splitting by port, instead of accepting both protocols on one port, means a client that sends the wrong protocol fails the connection outright instead of silently succeeding on the wrong protocol.
func NewServers(http1Addr, http2Addr string, cached *cache.Cache, shapeCached *cache.ShapeCache) (http1Server, http2Server *http.Server) {
	svc := service.New(cached)
	h := handler.New(svc)

	shapeSvc := service.NewShape(shapeCached)
	shapeH := handler.NewShape(shapeSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /rest/json/http1/order/depth-zero", h.ServeJSONDepthZero)
	mux.HandleFunc("GET /rest/json/http1/order/depth-two", h.ServeJSONDepthTwo)
	mux.HandleFunc("GET /rest/json/http1/order/depth-four", h.ServeJSONDepthFour)
	mux.HandleFunc("GET /rest/json/http1/order/one", h.ServeJSONOne)
	mux.HandleFunc("GET /rest/json/http1/order/hundred", h.ServeJSONHundred)
	mux.HandleFunc("GET /rest/json/http1/order/thousand", h.ServeJSONThousand)
	mux.HandleFunc("GET /rest/json/http2/order/depth-zero", h.ServeJSONDepthZero)
	mux.HandleFunc("GET /rest/json/http2/order/depth-two", h.ServeJSONDepthTwo)
	mux.HandleFunc("GET /rest/json/http2/order/depth-four", h.ServeJSONDepthFour)
	mux.HandleFunc("GET /rest/json/http2/order/one", h.ServeJSONOne)
	mux.HandleFunc("GET /rest/json/http2/order/hundred", h.ServeJSONHundred)
	mux.HandleFunc("GET /rest/json/http2/order/thousand", h.ServeJSONThousand)
	mux.HandleFunc("GET /rest/protobuf/http1/order/depth-zero", h.ServeProtobufDepthZero)
	mux.HandleFunc("GET /rest/protobuf/http1/order/depth-two", h.ServeProtobufDepthTwo)
	mux.HandleFunc("GET /rest/protobuf/http1/order/depth-four", h.ServeProtobufDepthFour)
	mux.HandleFunc("GET /rest/protobuf/http1/order/one", h.ServeProtobufOne)
	mux.HandleFunc("GET /rest/protobuf/http1/order/hundred", h.ServeProtobufHundred)
	mux.HandleFunc("GET /rest/protobuf/http1/order/thousand", h.ServeProtobufThousand)
	mux.HandleFunc("GET /rest/protobuf/http2/order/depth-zero", h.ServeProtobufDepthZero)
	mux.HandleFunc("GET /rest/protobuf/http2/order/depth-two", h.ServeProtobufDepthTwo)
	mux.HandleFunc("GET /rest/protobuf/http2/order/depth-four", h.ServeProtobufDepthFour)
	mux.HandleFunc("GET /rest/protobuf/http2/order/one", h.ServeProtobufOne)
	mux.HandleFunc("GET /rest/protobuf/http2/order/hundred", h.ServeProtobufHundred)
	mux.HandleFunc("GET /rest/protobuf/http2/order/thousand", h.ServeProtobufThousand)

	mux.HandleFunc("GET /rest/json/http1/shape/depth0/compact", shapeH.ServeJSONDepth0Compact)
	mux.HandleFunc("GET /rest/json/http1/shape/depth0/large", shapeH.ServeJSONDepth0Large)
	mux.HandleFunc("GET /rest/json/http1/shape/depth1-wide/compact", shapeH.ServeJSONDepth1WideCompact)
	mux.HandleFunc("GET /rest/json/http1/shape/depth1-wide/large", shapeH.ServeJSONDepth1WideLarge)
	mux.HandleFunc("GET /rest/json/http1/shape/depth3-narrow/compact", shapeH.ServeJSONDepth3NarrowCompact)
	mux.HandleFunc("GET /rest/json/http1/shape/depth3-narrow/large", shapeH.ServeJSONDepth3NarrowLarge)
	mux.HandleFunc("GET /rest/json/http1/shape/depth4-wide/compact", shapeH.ServeJSONDepth4WideCompact)
	mux.HandleFunc("GET /rest/json/http1/shape/depth4-wide/large", shapeH.ServeJSONDepth4WideLarge)
	mux.HandleFunc("GET /rest/json/http2/shape/depth0/compact", shapeH.ServeJSONDepth0Compact)
	mux.HandleFunc("GET /rest/json/http2/shape/depth0/large", shapeH.ServeJSONDepth0Large)
	mux.HandleFunc("GET /rest/json/http2/shape/depth1-wide/compact", shapeH.ServeJSONDepth1WideCompact)
	mux.HandleFunc("GET /rest/json/http2/shape/depth1-wide/large", shapeH.ServeJSONDepth1WideLarge)
	mux.HandleFunc("GET /rest/json/http2/shape/depth3-narrow/compact", shapeH.ServeJSONDepth3NarrowCompact)
	mux.HandleFunc("GET /rest/json/http2/shape/depth3-narrow/large", shapeH.ServeJSONDepth3NarrowLarge)
	mux.HandleFunc("GET /rest/json/http2/shape/depth4-wide/compact", shapeH.ServeJSONDepth4WideCompact)
	mux.HandleFunc("GET /rest/json/http2/shape/depth4-wide/large", shapeH.ServeJSONDepth4WideLarge)
	mux.HandleFunc("GET /rest/protobuf/http1/shape/depth0/compact", shapeH.ServeProtobufDepth0Compact)
	mux.HandleFunc("GET /rest/protobuf/http1/shape/depth0/large", shapeH.ServeProtobufDepth0Large)
	mux.HandleFunc("GET /rest/protobuf/http1/shape/depth1-wide/compact", shapeH.ServeProtobufDepth1WideCompact)
	mux.HandleFunc("GET /rest/protobuf/http1/shape/depth1-wide/large", shapeH.ServeProtobufDepth1WideLarge)
	mux.HandleFunc("GET /rest/protobuf/http1/shape/depth3-narrow/compact", shapeH.ServeProtobufDepth3NarrowCompact)
	mux.HandleFunc("GET /rest/protobuf/http1/shape/depth3-narrow/large", shapeH.ServeProtobufDepth3NarrowLarge)
	mux.HandleFunc("GET /rest/protobuf/http1/shape/depth4-wide/compact", shapeH.ServeProtobufDepth4WideCompact)
	mux.HandleFunc("GET /rest/protobuf/http1/shape/depth4-wide/large", shapeH.ServeProtobufDepth4WideLarge)
	mux.HandleFunc("GET /rest/protobuf/http2/shape/depth0/compact", shapeH.ServeProtobufDepth0Compact)
	mux.HandleFunc("GET /rest/protobuf/http2/shape/depth0/large", shapeH.ServeProtobufDepth0Large)
	mux.HandleFunc("GET /rest/protobuf/http2/shape/depth1-wide/compact", shapeH.ServeProtobufDepth1WideCompact)
	mux.HandleFunc("GET /rest/protobuf/http2/shape/depth1-wide/large", shapeH.ServeProtobufDepth1WideLarge)
	mux.HandleFunc("GET /rest/protobuf/http2/shape/depth3-narrow/compact", shapeH.ServeProtobufDepth3NarrowCompact)
	mux.HandleFunc("GET /rest/protobuf/http2/shape/depth3-narrow/large", shapeH.ServeProtobufDepth3NarrowLarge)
	mux.HandleFunc("GET /rest/protobuf/http2/shape/depth4-wide/compact", shapeH.ServeProtobufDepth4WideCompact)
	mux.HandleFunc("GET /rest/protobuf/http2/shape/depth4-wide/large", shapeH.ServeProtobufDepth4WideLarge)

	mux.HandleFunc("GET /rest/json/http1/shape/ping", shapeH.ServePingJSON)
	mux.HandleFunc("GET /rest/json/http2/shape/ping", shapeH.ServePingJSON)
	mux.HandleFunc("GET /rest/protobuf/http1/shape/ping", shapeH.ServePingProtobuf)
	mux.HandleFunc("GET /rest/protobuf/http2/shape/ping", shapeH.ServePingProtobuf)

	http1Server = &http.Server{
		Addr:    http1Addr,
		Handler: mux,
	}
	http1Server.Protocols = new(http.Protocols)
	http1Server.Protocols.SetHTTP1(true)

	http2Server = &http.Server{
		Addr:    http2Addr,
		Handler: mux,
	}
	http2Server.Protocols = new(http.Protocols)
	http2Server.Protocols.SetUnencryptedHTTP2(true)

	return http1Server, http2Server
}
