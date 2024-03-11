package interfaces

import (
	"context"
	"database/sql"
)

// UserScopeService defines the interface for user-scope related operations.
type UserScopeService interface {
	UpsertUserScope(ctx context.Context, userID, scopeID int64, role string, otx ...*sql.Tx) error
	GetUserScope(ctx context.Context, userID, scopeID int64, otx ...*sql.Tx) (*UserScope, error)
	ValidateUserScope(ctx context.Context, userID, scopeID int64, role string, otx ...*sql.Tx) bool
	GetUserScopesByRole(ctx context.Context, userID int64, role string, otx ...*sql.Tx) ([]UserScope, error)
	DeleteUserScope(ctx context.Context, userID, scopeID int64, otx ...*sql.Tx) error
}

// UserScope represents a user-scope relationship.
type UserScope struct {
	UserID  int64  `json:"user_id"`
	ScopeID int64  `json:"scope_id"`
	Role    string `json:"role"`
}
