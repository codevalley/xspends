package interfaces

import (
	"context"
	"database/sql"
	"time"
)

type Transaction struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	SourceID    int64     `json:"source_id"`
	Tags        []string  `json:"tags"`
	CategoryID  int64     `json:"category_id"`
	Timestamp   time.Time `json:"timestamp"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	ScopeID     int64     `json:"scope_id"`
}

type TransactionFilter struct {
	UserID       int64
	StartDate    string
	EndDate      string
	Tags         []string
	Category     string
	Type         string
	Description  string
	MinAmount    float64
	MaxAmount    float64
	SortBy       string
	SortOrder    string // "ASC" or "DESC"
	Page         int
	ItemsPerPage int
}

// TransactionService defines the interface for transaction operations.
type TransactionService interface {
	InsertTransaction(ctx context.Context, txn Transaction, otx ...*sql.Tx) error
	UpdateTransaction(ctx context.Context, txn Transaction, otx ...*sql.Tx) error
	DeleteTransaction(ctx context.Context, transactionID int64, userID int64, otx ...*sql.Tx) error
	GetTransactionByID(ctx context.Context, transactionID int64, userID int64, otx ...*sql.Tx) (*Transaction, error)
	GetTransactionsByFilter(ctx context.Context, filter TransactionFilter, otx ...*sql.Tx) ([]Transaction, error)
}
