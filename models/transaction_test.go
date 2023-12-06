package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// InsertTransaction successfully inserts a new transaction with valid input
func TestInsertTransactionSuccess(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	txn := Transaction{
		UserID:      1,
		SourceID:    1,
		CategoryID:  1,
		Amount:      100.0,
		Type:        "debit",
		Description: "Test transaction",
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := InsertTransaction(ctx, txn, mockDBService)
	assert.NoError(t, err)
}

// UpdateTransaction successfully updates an existing transaction with valid input
func TestUpdateTransactionSuccess(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	txn := Transaction{
		ID:          1,
		UserID:      1,
		SourceID:    1,
		CategoryID:  1,
		Amount:      200.0,
		Type:        "credit",
		Description: "Updated transaction",
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := UpdateTransaction(ctx, txn, mockDBService)
	assert.NoError(t, err)
}

// // InsertTransaction returns an error when snowflakeGenerator is nil
// func TestInsertTransactionErrorSnowflakeGeneratorNil(t *testing.T) {
// 	tearDown := setUp(t)
// 	defer tearDown()

// 	snowflakeGenerator = nil

// 	txn := Transaction{
// 		UserID:      1,
// 		SourceID:    1,
// 		CategoryID:  1,
// 		Amount:      100.0,
// 		Type:        "debit",
// 		Description: "Test transaction",
// 	}

// 	err := InsertTransaction(ctx, txn, mockDBService)
// 	assert.Error(t, err)
// 	assert.EqualError(t, err, "Snowflake generator is not initialized")
// }

// // UpdateTransaction returns an error when snowflakeGenerator is nil
// func TestUpdateTransactionErrorSnowflakeGeneratorNil(t *testing.T) {
// 	tearDown := setUp(t)
// 	defer tearDown()

// 	snowflakeGenerator = nil

// 	txn := Transaction{
// 		ID:          1,
// 		UserID:      1,
// 		SourceID:    1,
// 		CategoryID:  1,
// 		Amount:      200.0,
// 		Type:        "credit",
// 		Description: "Updated transaction",
// 	}

// 	err := UpdateTransaction(ctx, txn, mockDBService)
// 	assert.Error(t, err)
// 	assert.EqualError(t, err, "Snowflake generator is not initialized")
// }

// GetTransactionByID returns an error when the transaction does not exist
func TestGetTransactionByIDErrorTransactionNotExist(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	mockExecutor.EXPECT().
		QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, sql.ErrNoRows).
		Times(1)

	transaction, err := GetTransactionByID(ctx, 1, 1, mockDBService)
	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.EqualError(t, err, "get transaction by ID failed")
}

// GetTransactionsByFilter returns an error when the query fails
func TestGetTransactionsByFilterErrorQueryFailed(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	mockExecutor.EXPECT().
		QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("query error")).
		Times(1)

	filter := TransactionFilter{
		UserID: 1,
	}

	transactions, err := GetTransactionsByFilter(ctx, filter, mockDBService)
	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.EqualError(t, err, "querying transactions by filter failed")
}

// DeleteTransaction successfully deletes an existing transaction with valid input
// func TestDeleteTransaction_SuccessfullyDeletesExistingTransactionWithValidInput(t *testing.T) {
// 	tearDown := setUp(t)
// 	defer tearDown()

// 	// Mock getExecutor function
// 	getExecutor = func(dbService *DBService, otx ...*sql.Tx) (bool, DBExecutor) {
// 		return false, mockExecutor
// 	}

// 	// Create a test transaction
// 	transaction := Transaction{
// 		ID:          1,
// 		UserID:      1,
// 		SourceID:    1,
// 		CategoryID:  1,
// 		Timestamp:   time.Now(),
// 		Amount:      100.0,
// 		Type:        "expense",
// 		Description: "Test transaction",
// 	}

// 	// Call the DeleteTransaction function
// 	err := DeleteTransaction(ctx, transaction.ID, transaction.UserID, mockDBService)
// 	assert.NoError(t, err)
// }

// GetTransactionByID successfully retrieves an existing transaction with valid input
// func TestGetTransactionByIDSuccess(t *testing.T) {
// 	tearDown := setUp(t)
// 	defer tearDown()

// 	transaction := &Transaction{
// 		ID:          1,
// 		UserID:      1,
// 		SourceID:    1,
// 		CategoryID:  1,
// 		Timestamp:   time.Now(),
// 		Amount:      100.0,
// 		Type:        "debit",
// 		Description: "Test Transaction",
// 	}
// 	mockRow := sqlmock.NewRows([]string{"id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description"}).
// 	AddRow(1, 1, 1, 1, time.Now(), 10.0, "Expense", "Groceries")
// 	mockExecutor.EXPECT().
// 		QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
// 		Return(mockRow).
// 		Times(1)
// 	mockRow.EXPECT().
// 		Scan(gomock.Any()).
// 		SetArg(0, *transaction).
// 		Return(nil).
// 		Times(1)
// 	mockDBService.EXPECT().
// 		GetTagsByTransactionID(gomock.Any(), transaction.ID, gomock.Any()).
// 		Return([]Tag{}, nil).
// 		Times(1)

// 	result, err := GetTransactionByID(ctx, transaction.ID, transaction.UserID, mockDBService)
// 	assert.NoError(t, err)
// 	assert.Equal(t, transaction, result)
// }

// GetTransactionsByFilter successfully retrieves transactions with valid input
func TestGetTransactionsByFilterSuccess(t *testing.T) {
	db, mock, err := sqlmock.New() // mock sql.DB
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create a mock DBService using the mock database connection
	mockDBService := &DBService{Executor: db}

	filter := TransactionFilter{
		UserID:       490196676401692948,
		SortOrder:    SortOrderDesc,
		Page:         1,
		ItemsPerPage: 10,
	}

	mockRows := sqlmock.NewRows([]string{"id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description"}).
		AddRow(1, 490196676401692948, 1, 1, time.Now(), 10.0, "Expense", "Groceries")

	mock.ExpectQuery("SELECT id, user_id, source_id, category_id, timestamp, amount, type, description FROM transactions WHERE user_id =").
		WithArgs(filter.UserID).
		WillReturnRows(mockRows)

	// Call the function under test
	transactions, err := GetTransactionsByFilter(ctx, filter, mockDBService)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, int64(1), transactions[0].ID)
	assert.Equal(t, int64(490196676401692948), transactions[0].UserID)
	assert.Equal(t, int64(1), transactions[0].SourceID)
	assert.Equal(t, int64(1), transactions[0].CategoryID)
	assert.NotZero(t, transactions[0].Timestamp)
	assert.Equal(t, 10.0, transactions[0].Amount)
	assert.Equal(t, "Expense", transactions[0].Type)
	assert.Equal(t, "Groceries", transactions[0].Description)
}
