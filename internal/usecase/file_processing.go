package usecase

import (
	"context"
	"time"

	"github.com/bekzod003/image-grpc/internal/entity"
	"github.com/bekzod003/image-grpc/pkg/logger"
	ratelimiter "github.com/bekzod003/image-grpc/pkg/rate_limiter"
	"github.com/google/uuid"
)

type fileProcessingRepo interface {
	Store(ctx context.Context, file *entity.File) (*entity.File, error)
	Files(ctx context.Context, pagination *entity.LimitOffset) ([]*entity.File, error)
	Download(ctx context.Context, id string) (*entity.File, error)
	Count(ctx context.Context) (int32, error)
}

type fileProcessingUseCase struct {
	repo fileProcessingRepo
	log  logger.LoggerI
}

func NewFileProcessingUseCase(repo fileProcessingRepo, log logger.LoggerI) *fileProcessingUseCase {
	return &fileProcessingUseCase{repo: repo, log: log}
}

func (u *fileProcessingUseCase) Store(ctx context.Context, req *entity.StoreFileRequest) (*entity.File, error) {
	u.log.Info("Store file usecase been called", logger.Any("req", req))

	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// check store limit
	limitReached := ratelimiter.LimitDownloadAndStore(req.UserID)
	if limitReached {
		return nil, entity.ErrLimitReached
	}
	defer ratelimiter.LimitDownloadAndStoreDone(req.UserID)
	// finish check

	// store file
	file := &entity.File{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Link:      req.Link,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	file, err := u.repo.Store(c, file)
	if err != nil {
		u.log.Error("Store file usecase error", logger.Error(err))
		return nil, err
	}

	u.log.Info("Store file usecase success", logger.Any("file", file))
	return file, nil
}

func (u *fileProcessingUseCase) Files(ctx context.Context, req *entity.FilesRequest) ([]*entity.File, error) {
	u.log.Info("Files usecase been called", logger.Any("pagination", req))

	// check limit
	limitReached := ratelimiter.LimitList(req.UserID)
	if limitReached {
		return nil, entity.ErrLimitReached
	}
	defer ratelimiter.LimitListDone(req.UserID)
	//finish limit

	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	files, err := u.repo.Files(c, req.Pagination)
	if err != nil {
		u.log.Error("Files usecase error", logger.Error(err))
		return nil, err
	}

	u.log.Info("Files usecase success", logger.Any("files", files))
	return files, nil
}

func (u *fileProcessingUseCase) Download(ctx context.Context, req *entity.DownloadRequest) (*entity.File, error) {
	u.log.Info("Download usecase been called", logger.Any("req", req))

	// check limit
	limitReached := ratelimiter.LimitDownloadAndStore(req.UserID)
	if limitReached {
		return nil, entity.ErrLimitReached
	}
	defer ratelimiter.LimitDownloadAndStoreDone(req.UserID)
	//finish limit

	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	file, err := u.repo.Download(c, req.FileID)
	if err != nil {
		u.log.Error("Download usecase error", logger.Error(err))
		return nil, err
	}

	u.log.Info("Download usecase success", logger.Any("file", file))
	return file, nil
}
