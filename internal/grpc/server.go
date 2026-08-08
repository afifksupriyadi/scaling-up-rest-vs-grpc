package grpc

import (
	"google.golang.org/grpc"

	"scaling-up-rest-vs-grpc/internal/data/cache"
	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/grpc/handler"
	"scaling-up-rest-vs-grpc/internal/grpc/service"
)

// initialWindowSize sets the HTTP/2 flow control window to 4 MB, matching the practical ceiling recommended by the grpc-go maintainers, since most TCP connections cannot effectively utilize a larger window regardless.
// - Without this, grpc-go's dynamic BDP-based window sizing starts small and ramps up gradually, which never reaches its optimal size under the short-lived, high-concurrency request pattern used in this benchmark.
// - net/http's REST server does not need this adjustment, since it uses a large static window from the very first request instead of growing dynamically.
const initialWindowSize = 1024 * 1024 * 4 // 4 MB

// NewServer builds a *grpc.Server without TLS (plaintext), consistent with grpc-go's own default behavior when grpc.Creds() is not called.
func NewServer(cached *cache.Cache, shapeCached *cache.ShapeCache) *grpc.Server {
	svc := service.New(cached)
	h := handler.New(svc)

	shapeSvc := service.NewShape(shapeCached)
	shapeH := handler.NewShape(shapeSvc)

	s := grpc.NewServer(
		grpc.InitialWindowSize(initialWindowSize),
		grpc.InitialConnWindowSize(initialWindowSize),
	)
	model.RegisterStudentServiceServer(s, h)
	model.RegisterShapeExperimentServiceServer(s, shapeH)
	return s
}
