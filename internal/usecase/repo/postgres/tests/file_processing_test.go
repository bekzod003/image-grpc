package tests

import (
	"context"
	"testing"
	"time"

	"github.com/bekzod003/image-grpc/internal/entity"
	"github.com/bekzod003/image-grpc/internal/usecase/repo/postgres"
	"github.com/google/uuid"
)

func TestFileStore(t *testing.T) {
	repo := postgres.NewFileRepo(client)

	f, err := repo.Store(context.Background(), &entity.File{
		ID:        uuid.New().String(),
		Name:      "random name",
		Link:      "https://www.google.com/",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		t.Error("error storing file ->", err)
	}
	t.Logf("file -> \n%+v", f)
}

func TestFiles(t *testing.T) {
	repo := postgres.NewFileRepo(client)

	f, err := repo.Files(context.Background(), nil)
	if err != nil {
		t.Error("error getting files ->", err)
	}
	for _, file := range f {
		t.Logf("file -> \n%+v", file)
	}
}

func TestFileDownload(t *testing.T) {
	repo := postgres.NewFileRepo(client)

	f, err := repo.Download(context.Background(), "5e0b3846-66b4-474f-9622-361b6b770db0")
	if err != nil {
		t.Error("error downloading file ->", err)
	}
	t.Logf("file -> \n%+v", f)
}

func TestFileCount(t *testing.T) {
	repo := postgres.NewFileRepo(client)

	f, err := repo.Count(context.Background())
	if err != nil {
		t.Error("error getting file count ->", err)
	}
	t.Logf("file count -> \n%+v", f)
}
