package mock

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"

	"github.com/stretchr/testify/mock"
)

// MockTransactionTagModel is a mock implementation of the TransactionTagService interface for testing.
type MockTransactionTagModel struct {
	mock.Mock
}

// GetTagsByTransactionID mocks the GetTagsByTransactionID method.
func (m *MockTransactionTagModel) GetTagsByTransactionID(ctx context.Context, transactionID int64, otx ...*sql.Tx) ([]interfaces.Tag, error) {
	args := m.Called(ctx, transactionID, otx)
	return args.Get(0).([]interfaces.Tag), args.Error(1)
}

// InsertTransactionTag mocks the InsertTransactionTag method.
func (m *MockTransactionTagModel) InsertTransactionTag(ctx context.Context, transactionID, tagID int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, transactionID, tagID, otx)
	return args.Error(0)
}

// DeleteTransactionTag mocks the DeleteTransactionTag method.
func (m *MockTransactionTagModel) DeleteTransactionTag(ctx context.Context, transactionID, tagID int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, transactionID, tagID, otx)
	return args.Error(0)
}

// AddTagsToTransaction mocks the AddTagsToTransaction method.
func (m *MockTransactionTagModel) AddTagsToTransaction(ctx context.Context, transactionID int64, tags []string, scope []int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, transactionID, tags, scope, otx)
	return args.Error(0)
}

// UpdateTagsForTransaction mocks the UpdateTagsForTransaction method.
func (m *MockTransactionTagModel) UpdateTagsForTransaction(ctx context.Context, transactionID int64, tags []string, scope []int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, transactionID, tags, scope, otx)
	return args.Error(0)
}

// DeleteTagsFromTransaction mocks the DeleteTagsFromTransaction method.
func (m *MockTransactionTagModel) DeleteTagsFromTransaction(ctx context.Context, transactionID int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, transactionID, otx)
	return args.Error(0)
}

// Idiomatic interface compliance check.
// Ensure MockTransactionTagModel implements TransactionTagService
var _ interfaces.TransactionTagService = &MockTransactionTagModel{}
