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
