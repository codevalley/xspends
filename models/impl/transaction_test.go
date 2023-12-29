package impl

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"xspends/models/interfaces"
	xmock "xspends/models/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupNewMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred when creating a mock database connection: %v", err)
	}
	t.Cleanup(func() {
		db.Close() // ensure the db connection is closed after the test
	})
	mockModelService.DBService.Executor = db

	return db, mock
}
func setupForeignKeyMocks(mockM sqlmock.Sqlmock, txn interfaces.Transaction) {
	// Mock the user existence check
	mockM.ExpectQuery("^SELECT (.+) FROM users WHERE").
		WithArgs(txn.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow("1")) // indicating user exists

	// Mock the source existence check
	mockM.ExpectQuery(`SELECT 1 FROM sources WHERE id = \? AND user_id = \? LIMIT 1`).
		WithArgs(txn.SourceID, txn.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow("1")) // indicating source exists

	// Mock the category existence check
	mockM.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(txn.CategoryID, txn.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow("1")) // indicating category exists
}

func TestInsertTransactionV2(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	txn := interfaces.Transaction{
		UserID:      1,
		SourceID:    1,
		CategoryID:  1,
		Amount:      100.0,
		Type:        "expense",
		Description: "Groceries",
		Tags:        []string{"groceries", "food"},
	}

	db, mockM := setupNewMock(t)
	defer db.Close()

	t.Run("Successful Insertion", func(t *testing.T) {
		mockTagModel := new(xmock.MockTagModel)
		mockTransactionTagModel := new(xmock.MockTransactionTagModel)
		mockModelService.TransactionTagModel = mockTransactionTagModel
		mockModelService.TagModel = mockTagModel
		setupForeignKeyMocks(mockM, txn)
		mockM.ExpectExec("INSERT INTO transactions").
			WithArgs(sqlmock.AnyArg(), txn.UserID, txn.SourceID, txn.CategoryID, sqlmock.AnyArg(), txn.Amount, txn.Type, txn.Description).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Setup for tag handling
		for _, tag := range txn.Tags {
			mockTagModel.On(
				"GetTagByName",
				mock.Anything, tag, txn.UserID, mock.Anything,
			).Return(&interfaces.Tag{ID: 1, Name: tag, UserID: txn.UserID}, nil).Once()

			mockTagModel.On(
				"InsertTag",
				mock.Anything, mock.AnythingOfType("*interfaces.Tag"), mock.Anything,
			).Return(nil).Maybe()
		}

		mockTransactionTagModel.On(
			"AddTagsToTransaction",
			mock.Anything, mock.Anything, mock.Anything, txn.UserID, mock.Anything,
		).Return(nil).Once()

		// Call the method under test
		err := mockModelService.TransactionModel.InsertTransaction(context.Background(), txn)
		assert.NoError(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
		mockTagModel.AssertExpectations(t)
		mockTransactionTagModel.AssertExpectations(t)
	})
	// Subtest 2: Foreign Key Validation Failure (User does not exist)
	t.Run("Foreign Key Validation Failure - User Not Found", func(t *testing.T) {
		// Setup new mock database for clean expectation slate
		_, mock1 := setupNewMock(t)
		// Expect the user existence check query
		mock1.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(txn.UserID).WillReturnRows(sqlmock.NewRows([]string{"exists"}))

		err := mockModelService.TransactionModel.InsertTransaction(context.Background(), txn)

		// We are expecting an error because the user is not found
		assert.Error(t, err)

		// Validate that all expectations set on the mock were met
		assert.NoError(t, mock1.ExpectationsWereMet())
	})

	// Subtest 3: Insert Transaction Execution Failure
	t.Run("Insert Transaction Execution Failure", func(t *testing.T) {
		_, mockM = setupNewMock(t)
		setupForeignKeyMocks(mockM, txn) // Assumes a function to set up foreign key validation

		mockM.ExpectExec("INSERT INTO transactions").
			WithArgs(sqlmock.AnyArg(), txn.UserID, txn.SourceID, txn.CategoryID, sqlmock.AnyArg(), txn.Amount, txn.Type, txn.Description).
			WillReturnError(sql.ErrConnDone) // Simulate a connection error or similar

		err := mockModelService.TransactionModel.InsertTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	//Subtest 4: Handling Transaction Tags Fails
	t.Run("Handling Transaction Tags Fails", func(t *testing.T) {
		_, mockM = setupNewMock(t)
		mockTagModel := new(xmock.MockTagModel)
		mockTransactionTagModel := new(xmock.MockTransactionTagModel)
		mockModelService.TransactionTagModel = mockTransactionTagModel
		mockModelService.TagModel = mockTagModel
		setupForeignKeyMocks(mockM, txn) // Set up foreign key validations

		// Set up mocks for successful transaction insert
		mockM.ExpectExec("INSERT INTO transactions").
			WithArgs(sqlmock.AnyArg(), txn.UserID, txn.SourceID, txn.CategoryID, sqlmock.AnyArg(), txn.Amount, txn.Type, txn.Description).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Setup for tag handling
		for _, tag := range txn.Tags {
			mockTagModel.On(
				"GetTagByName",
				mock.Anything, tag, txn.UserID, mock.Anything,
			).Return(&interfaces.Tag{ID: 1, Name: tag, UserID: txn.UserID}, nil).Once()

			mockTagModel.On(
				"InsertTag",
				mock.Anything, mock.AnythingOfType("*interfaces.Tag"), mock.Anything,
			).Return(nil).Maybe()
		}
		// Simulate a failure in AddTagsToTransaction
		mockTransactionTagModel.On(
			"AddTagsToTransaction",
			mock.Anything, mock.Anything, mock.Anything, txn.UserID, mock.Anything,
		).Return(sql.ErrConnDone).Once() // Use an appropriate error

		err := mockModelService.TransactionModel.InsertTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
		mockTransactionTagModel.AssertExpectations(t)
	})
}

func TestUpdateTransactionV2(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	txn := interfaces.Transaction{
		ID:          1, // Assume an existing transaction ID for update
		UserID:      1,
		SourceID:    1,
		CategoryID:  1,
		Amount:      150.0,
		Type:        "income",
		Description: "Updated Groceries",
		Tags:        []string{"updatedTag1", "updatedTag2"},
	}

	db, mockM := setupNewMock(t)
	defer db.Close()

	t.Run("Successful Update", func(t *testing.T) {
		mockTagModel := new(xmock.MockTagModel)
		mockTransactionTagModel := new(xmock.MockTransactionTagModel)
		mockModelService.TransactionTagModel = mockTransactionTagModel
		mockModelService.TagModel = mockTagModel
		setupForeignKeyMocks(mockM, txn)

		// Set up mock for successful update
		mockM.ExpectExec("UPDATE transactions").
			WithArgs(txn.SourceID, txn.CategoryID, txn.Amount, txn.Type, txn.Description, txn.ID, txn.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Set up mocks for tag handling
		for _, tag := range txn.Tags {
			mockTagModel.On(
				"GetTagByName",
				mock.Anything, tag, txn.UserID, mock.Anything,
			).Return(&interfaces.Tag{ID: 1, Name: tag, UserID: txn.UserID}, nil).Once()

			mockTagModel.On(
				"InsertTag",
				mock.Anything, mock.AnythingOfType("*interfaces.Tag"), mock.Anything,
			).Return(nil).Maybe()
		}

		mockTransactionTagModel.On(
			"UpdateTagsForTransaction",
			mock.Anything, txn.ID, txn.Tags, txn.UserID, mock.Anything,
		).Return(nil).Once()

		// Call the method under test
		err := mockModelService.TransactionModel.UpdateTransaction(context.Background(), txn)
		assert.NoError(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
		mockTagModel.AssertExpectations(t)
		mockTransactionTagModel.AssertExpectations(t)
	})

	t.Run("Foreign Key Validation Failure - User Not Found", func(t *testing.T) {
		_, mockM = setupNewMock(t)

		// Set up user not found scenario
		mockM.ExpectQuery("^SELECT (.+) FROM users WHERE").
			WithArgs(txn.UserID).
			WillReturnRows(sqlmock.NewRows(nil)) // No rows returned to simulate user not found

		// Call the update method and expect an error
		err := mockModelService.TransactionModel.UpdateTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Update Transaction Execution Failure", func(t *testing.T) {
		_, mockM = setupNewMock(t)
		setupForeignKeyMocks(mockM, txn) // Assumes a function to set up foreign key validation

		// Simulate a failure during the execution of the update query
		mockM.ExpectExec("UPDATE transactions").
			WithArgs(txn.SourceID, txn.CategoryID, txn.Amount, txn.Type, txn.Description, txn.ID, txn.UserID).
			WillReturnError(sql.ErrConnDone)

		// Call the update method and expect an error
		err := mockModelService.TransactionModel.UpdateTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Handling Transaction Tags Update Fails", func(t *testing.T) {
		_, mockM = setupNewMock(t)
		mockTagModel := new(xmock.MockTagModel)
		mockTransactionTagModel := new(xmock.MockTransactionTagModel)
		mockModelService.TransactionTagModel = mockTransactionTagModel
		mockModelService.TagModel = mockTagModel
		setupForeignKeyMocks(mockM, txn) // Set up foreign key validations

		// Set up mocks for successful transaction update
		mockM.ExpectExec("UPDATE transactions").
			WithArgs(txn.SourceID, txn.CategoryID, txn.Amount, txn.Type, txn.Description, txn.ID, txn.UserID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Setup for tag handling
		for _, tag := range txn.Tags {
			mockTagModel.On(
				"GetTagByName",
				mock.Anything, tag, txn.UserID, mock.Anything,
			).Return(&interfaces.Tag{ID: 1, Name: tag, UserID: txn.UserID}, nil).Once()

			mockTagModel.On(
				"InsertTag",
				mock.Anything, mock.AnythingOfType("*interfaces.Tag"), mock.Anything,
			).Return(nil).Maybe()
		}

		// Simulate a failure in UpdateTagsForTransaction
		mockTransactionTagModel.On(
			"UpdateTagsForTransaction",
			mock.Anything, txn.ID, txn.Tags, txn.UserID, mock.Anything,
		).Return(sql.ErrConnDone).Once() // Use an appropriate error

		// Call the update method and expect an error
		err := mockModelService.TransactionModel.UpdateTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
		mockTransactionTagModel.AssertExpectations(t)
	})
}

func TestDeleteTransactionV2(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	transactionID := int64(1) // Assume an existing transaction ID for deletion
	userID := int64(1)

	db, mockM := setupNewMock(t)
	defer db.Close()

	t.Run("Successful Deletion", func(t *testing.T) {
		// Set up mock for successful deletion
		mockM.ExpectExec("DELETE FROM transactions").
			WithArgs(transactionID, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Call the method under test
		err := mockModelService.TransactionModel.DeleteTransaction(context.Background(), transactionID, userID)
		assert.NoError(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Delete Transaction Execution Failure", func(t *testing.T) {
		_, mockM = setupNewMock(t)

		// Simulate a failure during the execution of the delete query
		mockM.ExpectExec("DELETE FROM transactions").
			WithArgs(transactionID, userID).
			WillReturnError(sql.ErrConnDone)

		// Call the delete method and expect an error
		err := mockModelService.TransactionModel.DeleteTransaction(context.Background(), transactionID, userID)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Delete Non-Existent Transaction", func(t *testing.T) {
		_, mockM = setupNewMock(t)

		// Set up mock for a deletion attempt on a non-existent transaction
		mockM.ExpectExec("DELETE FROM transactions").
			WithArgs(transactionID, userID).
			WillReturnResult(sqlmock.NewResult(0, 0)) // Indicate no rows were affected

		// Call the delete method and check if error or some indication of non-existence is handled
		err := mockModelService.TransactionModel.DeleteTransaction(context.Background(), transactionID, userID)
		// Assert based on how your application should behave (error or just a no-op)
		// Here it is assumed the method won't error out if no transaction was found
		assert.NoError(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	// Add more subtests for different edge cases if necessary
	// Example: trying to delete with wrong user ID, etc.
}

func TestGetTransactionByIDV2(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	transactionID := int64(1) // Assume some transaction ID
	userID := int64(1)        // Assume some user ID

	// Assume a mock transaction to be returned
	mockTransaction := interfaces.Transaction{
		ID:          transactionID,
		UserID:      userID,
		SourceID:    1,
		CategoryID:  1,
		Timestamp:   time.Now(),
		Amount:      100.0,
		Type:        "expense",
		Description: "Mock Transaction",
	}

	db, mockM := setupNewMock(t)
	defer db.Close()

	t.Run("Successful Retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description"}).
			AddRow(mockTransaction.ID, mockTransaction.UserID, mockTransaction.SourceID, mockTransaction.CategoryID, mockTransaction.Timestamp, mockTransaction.Amount, mockTransaction.Type, mockTransaction.Description)

		mockM.ExpectQuery("SELECT (.+) FROM transactions WHERE").WithArgs(transactionID, userID).WillReturnRows(rows)

		// Assume getTagsForTransaction is properly implemented or mocked
		// Mock the call to getTagsForTransaction if it makes an external call

		transaction, err := mockModelService.TransactionModel.GetTransactionByID(context.Background(), transactionID, userID)
		assert.NoError(t, err)
		assert.Equal(t, mockTransaction, *transaction)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Transaction Not Found", func(t *testing.T) {
		mockM.ExpectQuery("SELECT (.+) FROM transactions WHERE").WithArgs(transactionID, userID).WillReturnRows(sqlmock.NewRows(nil)) // Return empty result set

		transaction, err := mockModelService.TransactionModel.GetTransactionByID(context.Background(), transactionID, userID)
		assert.Error(t, err) // Assuming function returns error on not found
		assert.Nil(t, transaction)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Query Execution Error", func(t *testing.T) {
		mockM.ExpectQuery("SELECT (.+) FROM transactions WHERE").WithArgs(transactionID, userID).WillReturnError(sql.ErrConnDone) // Simulate query execution error

		transaction, err := mockModelService.TransactionModel.GetTransactionByID(context.Background(), transactionID, userID)
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	// Add more subtests if needed, for example to cover the logic within getTagsForTransaction, etc.
}

func TestGetTransactionsByFilterV2(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	userID := int64(1) // Assume some user ID
	filter := interfaces.TransactionFilter{
		UserID: userID,
		// Set additional filter criteria as needed for testing
	}

	db, mockM := setupNewMock(t)
	defer db.Close()

	t.Run("Successful Retrieval", func(t *testing.T) {
		// Assume mock transactions to be returned
		mockTransactions := []interfaces.Transaction{
			// Define one or more mock transactions as per filter criteria
		}

		rows := sqlmock.NewRows([]string{"id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description"})
		for _, txn := range mockTransactions {
			rows = rows.AddRow(txn.ID, txn.UserID, txn.SourceID, txn.CategoryID, txn.Timestamp, txn.Amount, txn.Type, txn.Description)
		}

		mockM.ExpectQuery("SELECT (.+) FROM transactions").WillReturnRows(rows)

		// Assume getTagsForTransaction is properly implemented or mocked

		transactions, err := mockModelService.TransactionModel.GetTransactionsByFilter(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, mockTransactions, transactions)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Query Execution Error", func(t *testing.T) {
		mockM.ExpectQuery("SELECT (.+) FROM transactions").WillReturnError(sql.ErrConnDone) // Simulate query execution error

		transactions, err := mockModelService.TransactionModel.GetTransactionsByFilter(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, transactions)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Row Scan Error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description"}).AddRow(1, userID, 1, 1, time.Now(), 100.0, "expense", "Description")
		mockM.ExpectQuery("SELECT (.+) FROM transactions").WillReturnRows(rows)

		_ = rows.RowError(0, sql.ErrConnDone) // Simulate row scan error on the first row
		transactions, err := mockModelService.TransactionModel.GetTransactionsByFilter(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, transactions)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	// Add more subtests if needed, for example to cover different filter criteria and edge cases
}
