package grpc

import (
	"google.golang.org/grpc"

	"scaling-up-rest-vs-grpc/internal/data/cache"
	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/grpc/handler"
	"scaling-up-rest-vs-grpc/internal/grpc/service"
)

// NewServer builds a *grpc.Server without TLS (plaintext), consistent
// with grpc-go's own default behavior when grpc.Creds() is not called.
func NewServer(cached *cache.Cache) *grpc.Server {
	svc := service.New(cached)
	h := handler.New(svc)

	s := grpc.NewServer()
	model.RegisterStudentServiceServer(s, h)
	return s
}
