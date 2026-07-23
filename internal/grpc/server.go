package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"scaling-up-rest-vs-grpc/internal/data/model"
	"scaling-up-rest-vs-grpc/internal/data/store"
	"scaling-up-rest-vs-grpc/internal/grpc/handler"
	"scaling-up-rest-vs-grpc/internal/grpc/service"
)

// NewServer builds a *grpc.Server with TLS enabled and StudentService
// registered on it.
func NewServer(certPath, keyPath string, dataStore *store.Store) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	svc := service.New(dataStore)
	h := handler.New(svc)

	s := grpc.NewServer(grpc.Creds(creds))
	model.RegisterStudentServiceServer(s, h)
	return s, nil
}
