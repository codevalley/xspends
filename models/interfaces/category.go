package interfaces

import (
	"context"
	"database/sql"
	"time"
)

type Category struct {
	ID          int64     `json:"category_id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	ScopeID     int64     `json:"scope_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryService interface {
	InsertCategory(ctx context.Context, category *Category, otx ...*sql.Tx) error
	UpdateCategory(ctx context.Context, category *Category, otx ...*sql.Tx) error
	DeleteCategory(ctx context.Context, categoryID int64, scopes []int64, otx ...*sql.Tx) error
	GetCategoryByID(ctx context.Context, categoryID int64, scopes []int64, otx ...*sql.Tx) (*Category, error)
	GetScopedCategories(ctx context.Context, page int, itemsPerPage int, scopes []int64, otx ...*sql.Tx) ([]Category, error)
	CategoryIDExists(ctx context.Context, categoryID int64, scopes []int64, otx ...*sql.Tx) (bool, error)
}
