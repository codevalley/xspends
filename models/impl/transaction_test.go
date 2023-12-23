package impl

import (
	"context"
	"database/sql"
	"testing"
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
