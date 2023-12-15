package interfaces

import (
	"context"
	"database/sql"
	"time"
)

type Tag struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginationParams struct {
	Limit  int
	Offset int
}

type TagService interface {
	InsertTag(ctx context.Context, tag *Tag, otx ...*sql.Tx) error
	UpdateTag(ctx context.Context, tag *Tag, otx ...*sql.Tx) error
	DeleteTag(ctx context.Context, tagID int64, userID int64, otx ...*sql.Tx) error
	GetTagByID(ctx context.Context, tagID int64, userID int64, otx ...*sql.Tx) (*Tag, error)
	GetAllTags(ctx context.Context, userID int64, pagination PaginationParams, otx ...*sql.Tx) ([]Tag, error)
	GetTagByName(ctx context.Context, name string, userID int64, otx ...*sql.Tx) (*Tag, error)
}
