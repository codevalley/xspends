package mock

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"

	"github.com/stretchr/testify/mock"
)

// MockTransactionModel is a mock implementation of the TransactionService interface.
type MockTransactionModel struct {
	mock.Mock
}

// Ensure MockTransactionModel implements TransactionService.
var _ interfaces.TransactionService = &MockTransactionModel{}

func (m *MockTransactionModel) InsertTransaction(ctx context.Context, txn interfaces.Transaction, otx ...*sql.Tx) error {
	args := m.Called(ctx, txn, otx)
	return args.Error(0)
}

func (m *MockTransactionModel) UpdateTransaction(ctx context.Context, txn interfaces.Transaction, otx ...*sql.Tx) error {
	args := m.Called(ctx, txn, otx)
	return args.Error(0)
}

func (m *MockTransactionModel) DeleteTransaction(ctx context.Context, transactionID int64, userID int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, transactionID, userID, otx)
	return args.Error(0)
}

func (m *MockTransactionModel) GetTransactionByID(ctx context.Context, transactionID int64, userID int64, otx ...*sql.Tx) (*interfaces.Transaction, error) {
	args := m.Called(ctx, transactionID, userID, otx)
	return args.Get(0).(*interfaces.Transaction), args.Error(1)
}

func (m *MockTransactionModel) GetTransactionsByFilter(ctx context.Context, filter interfaces.TransactionFilter, otx ...*sql.Tx) ([]interfaces.Transaction, error) {
	args := m.Called(ctx, filter, otx)
	return args.Get(0).([]interfaces.Transaction), args.Error(1)
}

func (m *MockTransactionModel) InsertTransactionNew(ctx context.Context, txn interfaces.Transaction, otx ...*sql.Tx) error {
	args := m.Called(ctx, txn, otx)
	return args.Error(0)
}

func (m *MockTransactionModel) UpdateTransactionNew(ctx context.Context, txn interfaces.Transaction, otx ...*sql.Tx) error {
	args := m.Called(ctx, txn, otx)
	return args.Error(0)
}

func (m *MockTransactionModel) DeleteTransactionNew(ctx context.Context, transactionID int64, scopes []int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, transactionID, scopes, otx)
	return args.Error(0)
}

func (m *MockTransactionModel) GetTransactionByIDNew(ctx context.Context, transactionID int64, scopes []int64, otx ...*sql.Tx) (*interfaces.Transaction, error) {
	args := m.Called(ctx, transactionID, scopes, otx)
	return args.Get(0).(*interfaces.Transaction), args.Error(1)
}
