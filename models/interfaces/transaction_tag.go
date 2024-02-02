package interfaces

import (
	"context"
	"database/sql"
	"time"
)

// mostly obselete, this struct is not used anywhere.
type TransactionTag struct {
	TransactionID int64     `json:"transaction_id"`
	TagID         int64     `json:"tag_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TransactionTagService defines the interface for operations on transaction tags.
type TransactionTagService interface {
	GetTagsByTransactionID(ctx context.Context, transactionID int64, otx ...*sql.Tx) ([]Tag, error)
	InsertTransactionTag(ctx context.Context, transactionID, tagID int64, otx ...*sql.Tx) error
	DeleteTransactionTag(ctx context.Context, transactionID, tagID int64, otx ...*sql.Tx) error
	DeleteTagsFromTransaction(ctx context.Context, transactionID int64, otx ...*sql.Tx) error
	//deprecated methods
	AddTagsToTransaction(ctx context.Context, transactionID int64, tags []string, userID int64, otx ...*sql.Tx) error
	UpdateTagsForTransaction(ctx context.Context, transactionID int64, tags []string, userID int64, otx ...*sql.Tx) error
	//deprecated methods end
	AddTagsToTransactionNew(ctx context.Context, transactionID int64, tags []string, scopes []int64, otx ...*sql.Tx) error
	UpdateTagsForTransactionNew(ctx context.Context, transactionID int64, tags []string, scopes []int64, otx ...*sql.Tx) error
}
