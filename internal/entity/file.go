package entity

import (
	"time"
)

type File struct {
	ID   string
	Name string
	Link string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type StoreFileRequest struct {
	UserID string

	Name string
	Link string
}

type FilesRequest struct {
	Pagination *LimitOffset

	UserID string
}

type DownloadRequest struct {
	UserID string
	FileID string
}
