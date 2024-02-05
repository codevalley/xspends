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
	ScopeID   int64     `json:"scope_id"`
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

	DeleteTagNew(ctx context.Context, tagID int64, scopes []int64, otx ...*sql.Tx) error
	GetTagByIDNew(ctx context.Context, tagID int64, scopes []int64, otx ...*sql.Tx) (*Tag, error)
	GetScopedTags(ctx context.Context, scopes []int64, pagination PaginationParams, otx ...*sql.Tx) ([]Tag, error)
	GetTagByNameNew(ctx context.Context, name string, scopes []int64, otx ...*sql.Tx) (*Tag, error)
}
