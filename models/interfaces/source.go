package interfaces

import (
	"context"
	"database/sql"
	"time"
)

// Source struct as defined in your implementation.
type Source struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   float64   `json:"balance"`
	ScopeID   int64     `json:"scope_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SourceService defines the interface for source operations.
type SourceService interface {
	InsertSource(ctx context.Context, source *Source, otx ...*sql.Tx) error
	UpdateSource(ctx context.Context, source *Source, otx ...*sql.Tx) error

	DeleteSourceNew(ctx context.Context, sourceID int64, scopes []int64, otx ...*sql.Tx) error
	GetSourceByIDNew(ctx context.Context, sourceID int64, scopes []int64, otx ...*sql.Tx) (*Source, error)
	GetScopedSources(ctx context.Context, page int, itemsPerPage int, scopes []int64, otx ...*sql.Tx) ([]Source, error)
	GetSourcesNew(ctx context.Context, scopes []int64, otx ...*sql.Tx) ([]Source, error)
	SourceIDExistsNew(ctx context.Context, sourceID int64, scopes []int64, otx ...*sql.Tx) (bool, error)
}
