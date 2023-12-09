package interfaces

import (
	"context"
	"database/sql"
	"time"
)

type Category struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryService interface {
	InsertCategory(ctx context.Context, category *Category, otx ...*sql.Tx) error
	UpdateCategory(ctx context.Context, category *Category, otx ...*sql.Tx) error
	DeleteCategory(ctx context.Context, categoryID int64, userID int64, otx ...*sql.Tx) error
	GetAllCategories(ctx context.Context, userID int64, otx ...*sql.Tx) ([]Category, error)
	GetCategoryByID(ctx context.Context, categoryID int64, userID int64, otx ...*sql.Tx) (*Category, error)
	GetPagedCategories(ctx context.Context, page int, itemsPerPage int, userID int64, otx ...*sql.Tx) ([]Category, error)
	CategoryIDExists(ctx context.Context, categoryID int64, userID int64, otx ...*sql.Tx) (bool, error)
}
