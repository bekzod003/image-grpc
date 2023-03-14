package service

import (
	"context"

	"github.com/bekzod003/image-grpc/genproto/file_processing"
	"github.com/bekzod003/image-grpc/internal/entity"
	"github.com/bekzod003/image-grpc/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fileProcessingUsecase interface {
	Store(ctx context.Context, req *entity.StoreFileRequest) (*entity.File, error)
	Files(ctx context.Context, req *entity.FilesRequest) ([]*entity.File, error)
	Download(ctx context.Context, req *entity.DownloadRequest) (*entity.File, error)
}

type fileProcessingService struct {
	log                   logger.LoggerI
	fileProcessingUsecase fileProcessingUsecase
	file_processing.UnimplementedFileProcessingServer
}

// implementation of the service
func NewFileProcessingService(logger logger.LoggerI, usecase fileProcessingUsecase) file_processing.FileProcessingServer {
	return &fileProcessingService{
		log:                   logger,
		fileProcessingUsecase: usecase,
	}
}

func (s *fileProcessingService) Store(ctx context.Context, req *file_processing.StoreRequest) (*file_processing.File, error) {
	fileReq := &entity.StoreFileRequest{
		UserID: req.GetUserID(),
		Name:   req.GetName(),
		Link:   req.GetLink(),
	}

	file, err := s.fileProcessingUsecase.Store(ctx, fileReq)
	if err != nil {
		s.log.Error("error while storing file in grpc", logger.Error(err))
		return nil, err
	}

	return &file_processing.File{
		ID:        file.ID,
		Name:      file.Name,
		Link:      file.Link,
		CreatedAt: timestamppb.New(file.CreatedAt),
		UpdatedAt: timestamppb.New(file.UpdatedAt),
	}, nil
}

func (s *fileProcessingService) List(ctx context.Context, req *file_processing.ListRequest) (*file_processing.ListResponse, error) {
	listReq := &entity.FilesRequest{}
	files, err := s.fileProcessingUsecase.Files(ctx, listReq)
	if err != nil {
		s.log.Error("Error while getting files", logger.Error(err))
		return nil, err
	}

	filesResp := make([]*file_processing.File, len(files))
	for i, f := range files {
		filesResp[i] = &file_processing.File{
			ID:        f.ID,
			Name:      f.Name,
			Link:      f.Link,
			CreatedAt: timestamppb.New(f.CreatedAt),
			UpdatedAt: timestamppb.New(f.UpdatedAt),
		}
	}
	return &file_processing.ListResponse{Files: filesResp}, nil
}

func (s *fileProcessingService) Download(ctx context.Context, req *file_processing.DownloadRequest) (*file_processing.File, error) {
	downloadReq := &entity.DownloadRequest{
		UserID: req.GetUserID(),
		FileID: req.GetFileID(),
	}
	file, err := s.fileProcessingUsecase.Download(ctx, downloadReq)
	if err != nil {
		s.log.Error("Error while downloading file", logger.Error(err))
		return nil, err
	}

	return &file_processing.File{
		ID:        file.ID,
		Name:      file.Name,
		Link:      file.Link,
		CreatedAt: timestamppb.New(file.CreatedAt),
		UpdatedAt: timestamppb.New(file.UpdatedAt),
	}, nil
}
