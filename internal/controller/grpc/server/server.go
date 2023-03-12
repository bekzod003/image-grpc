package server

import (
	"github.com/bekzod003/image-grpc/config"
	"github.com/bekzod003/image-grpc/pkg/logger"
	"google.golang.org/grpc"
)

func NewGRPCServer(cfg *config.Config, log *logger.LoggerI) *grpc.Server {
	return grpc.NewServer(
	// Middleware
	// grpc.ChainUnaryInterceptor(),
	)
}
