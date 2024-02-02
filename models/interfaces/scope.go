package interfaces

import (
	"context"
	"database/sql"
)

type Scope struct {
	ScopeID int64  `json:"scope_id"`
	Type    string `json:"type"`
}

type ScopeService interface {
	CreateScope(ctx context.Context, scopeType string, otx ...*sql.Tx) (int64, error)
	GetScope(ctx context.Context, scopeID int64, otx ...*sql.Tx) (*Scope, error)
	DeleteScope(ctx context.Context, scopeID int64, otx ...*sql.Tx) error
	ScopeIDExists(ctx context.Context, scopeID int64, otx ...*sql.Tx) (bool, error)
}
