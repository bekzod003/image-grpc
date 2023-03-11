package postgres

import (
	"context"
	"fmt"

	"github.com/bekzod003/image-grpc/internal/entity"
	"github.com/bekzod003/image-grpc/pkg/database/client/postgresql"
)

type file struct {
	db postgresql.Client
}

func NewFileRepo(db postgresql.Client) *file {
	return &file{db: db}
}

// Store stores a file in the database
func (r *file) Store(ctx context.Context, file *entity.File) (*entity.File, error) {
	const query = `INSERT INTO files (id, name, link, created_at, updated_at) VALUES($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		file.ID, file.Name, file.Link, file.CreatedAt.UTC(), file.UpdatedAt.UTC(),
	)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// Files returns a list of files without URL. Pagination is supported.
// If pagination is nil, then all files are returned.
// File with link can be get using Download function.
func (r *file) Files(ctx context.Context, pagination *entity.LimitOffset) ([]*entity.File, error) {
	var (
		argsCount int
		args      []interface{}
	)

	query := `SELECT id, name, created_at, updated_at FROM files ORDER BY created_at DESC`
	if pagination != nil {
		if pagination.Limit != 0 {
			query += fmt.Sprintf(` LIMIT $%d`, argsCount+1)
			argsCount++
			args = append(args, pagination.Limit)
		}
		if pagination.Offset != 0 {
			query += fmt.Sprintf(` OFFSET $%d`, argsCount+1)
			argsCount++
			args = append(args, pagination.Offset)
		}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil || rows.Err() != nil {
		if rows.Err() == nil {
			err = rows.Err()
		}
		return nil, fmt.Errorf("error getting files -> %w", err)
	}
	defer rows.Close()

	var files []*entity.File
	for rows.Next() {
		var f entity.File
		if err := rows.Scan(&f.ID, &f.Name, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning file -> %w", err)
		}
		files = append(files, &f)
	}

	return files, nil
}

// Download returns a file with URL.
func (r *file) Download(ctx context.Context, id string) (*entity.File, error) {
	const query = `SELECT id, name, link, created_at, updated_at FROM files WHERE id = $1`
	var f entity.File
	if err := r.db.QueryRow(ctx, query, id).Scan(&f.ID, &f.Name, &f.Link, &f.CreatedAt, &f.UpdatedAt); err != nil {
		return nil, fmt.Errorf("error getting file -> %w", err)
	}
	return &f, nil
}

// Count returns the number of files in the database.
func (r *file) Count(ctx context.Context) (int32, error) {
	const query = `SELECT COUNT(*) FROM files`
	var count int32
	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("error getting files count -> %w", err)
	}
	return count, nil
}
