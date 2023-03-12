package service

import (
	"github.com/bekzod003/image-grpc/genproto/file_processing"
	"github.com/bekzod003/image-grpc/pkg/logger"
)

type fileProcessingService struct {
	logger.LoggerI
	// usecase
	file_processing.UnimplementedFileProcessingServer
}

// implementation of the service
func NewFileProcessingService(logger logger.LoggerI) file_processing.FileProcessingServer {
	return &fileProcessingService{
		LoggerI: logger,
	}
}
