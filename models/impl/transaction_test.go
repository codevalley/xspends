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
	ModelsService.DBService.Executor = db

	return db, mock
}
func setupForeignKeyMocks(mockM sqlmock.Sqlmock, txn interfaces.Transaction) {
	// Mock the user existence check
	mockM.ExpectQuery("^SELECT (.+) FROM users WHERE").
		WithArgs(txn.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow("1")) // indicating user exists

	// Mock the source existence check
	mockM.ExpectQuery("^SELECT (.+) FROM sources WHERE").
		WithArgs(txn.SourceID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow("1")) // indicating source exists

	// Mock the category existence check
	mockM.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(txn.CategoryID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow("1")) // indicating category exists

	// Mock the scope existence check
	mockM.ExpectQuery("^SELECT (.+) FROM scopes WHERE").
		WithArgs(txn.ScopeID).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow("1")) // indicating category exists
}

func TestInsertTransactionV2(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionModel = NewTransactionModel()
		config.UserModel = NewUserModel()
		config.SourceModel = NewSourceModel()
		config.CategoryModel = NewCategoryModel()
		config.ScopeModel = NewScopeModel()
	})
	defer tearDown()

	txn := interfaces.Transaction{
		UserID:      1,
		SourceID:    1,
		ScopeID:     1,
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
		ModelsService.TransactionTagModel = mockTransactionTagModel
		ModelsService.TagModel = mockTagModel
		setupForeignKeyMocks(mockM, txn)
		mockM.ExpectExec("INSERT INTO transactions").
			WithArgs(sqlmock.AnyArg(), txn.UserID, txn.SourceID, txn.CategoryID, sqlmock.AnyArg(), txn.Amount, txn.Type, txn.Description, txn.ScopeID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Setup for tag handling
		for _, tag := range txn.Tags {
			mockTagModel.On(
				"GetTagByName",
				mock.Anything, tag, []int64{txn.ScopeID}, mock.Anything,
			).Return(&interfaces.Tag{ID: 1, Name: tag, UserID: txn.UserID, ScopeID: txn.ScopeID}, nil).Once()

			mockTagModel.On(
				"InsertTag",
				mock.Anything, mock.AnythingOfType("*interfaces.Tag"), mock.Anything,
			).Return(nil).Maybe()
		}

		mockTransactionTagModel.On(
			"AddTagsToTransaction",
			mock.Anything, mock.Anything, mock.Anything, []int64{txn.ScopeID}, mock.Anything,
		).Return(nil).Once()

		// Call the method under test
		err := ModelsService.TransactionModel.InsertTransaction(context.Background(), txn)
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

		err := ModelsService.TransactionModel.InsertTransaction(context.Background(), txn)

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
			WithArgs(sqlmock.AnyArg(), txn.UserID, txn.SourceID, txn.CategoryID, sqlmock.AnyArg(), txn.Amount, txn.Type, txn.Description, txn.ScopeID).
			WillReturnError(sql.ErrConnDone) // Simulate a connection error or similar

		err := ModelsService.TransactionModel.InsertTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	//Subtest 4: Handling Transaction Tags Fails
	t.Run("Handling Transaction Tags Fails", func(t *testing.T) {
		_, mockM = setupNewMock(t)
		mockTagModel := new(xmock.MockTagModel)
		mockTransactionTagModel := new(xmock.MockTransactionTagModel)
		ModelsService.TransactionTagModel = mockTransactionTagModel
		ModelsService.TagModel = mockTagModel
		setupForeignKeyMocks(mockM, txn) // Set up foreign key validations

		// Set up mocks for successful transaction insert
		mockM.ExpectExec("INSERT INTO transactions").
			WithArgs(sqlmock.AnyArg(), txn.UserID, txn.SourceID, txn.CategoryID, sqlmock.AnyArg(), txn.Amount, txn.Type, txn.Description, txn.ScopeID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Setup for tag handling
		for _, tag := range txn.Tags {
			mockTagModel.On(
				"GetTagByName",
				mock.Anything, tag, []int64{txn.ScopeID}, mock.Anything,
			).Return(&interfaces.Tag{ID: 1, Name: tag, UserID: txn.UserID, ScopeID: txn.ScopeID}, nil).Once()

			mockTagModel.On(
				"InsertTag",
				mock.Anything, mock.AnythingOfType("*interfaces.Tag"), mock.Anything,
			).Return(nil).Maybe()
		}
		// Simulate a failure in AddTagsToTransaction
		mockTransactionTagModel.On(
			"AddTagsToTransaction",
			mock.Anything, mock.Anything, mock.Anything, []int64{txn.ScopeID}, mock.Anything,
		).Return(sql.ErrConnDone).Once() // Use an appropriate error

		err := ModelsService.TransactionModel.InsertTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
		mockTransactionTagModel.AssertExpectations(t)
	})
}

func TestUpdateTransactionV2(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionModel = NewTransactionModel()
		config.UserModel = NewUserModel()
		config.SourceModel = NewSourceModel()
		config.CategoryModel = NewCategoryModel()
		config.ScopeModel = NewScopeModel()
	})
	defer tearDown()

	txn := interfaces.Transaction{
		ID:          1, // Assume an existing transaction ID for update
		UserID:      1,
		SourceID:    1,
		ScopeID:     1,
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
		ModelsService.TransactionTagModel = mockTransactionTagModel
		ModelsService.TagModel = mockTagModel
		setupForeignKeyMocks(mockM, txn)

		// Set up mock for successful update
		mockM.ExpectExec("UPDATE transactions").
			WithArgs(txn.SourceID, txn.CategoryID, txn.Amount, txn.Type, txn.Description, txn.ID, txn.ScopeID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Set up mocks for tag handling
		for _, tag := range txn.Tags {
			mockTagModel.On(
				"GetTagByName",
				mock.Anything, tag, []int64{txn.ScopeID}, mock.Anything,
			).Return(&interfaces.Tag{ID: 1, Name: tag, UserID: txn.UserID, ScopeID: txn.ScopeID}, nil).Once()

			mockTagModel.On(
				"InsertTag",
				mock.Anything, mock.AnythingOfType("*interfaces.Tag"), mock.Anything,
			).Return(nil).Maybe()
		}

		mockTransactionTagModel.On(
			"UpdateTagsForTransaction",
			mock.Anything, txn.ID, txn.Tags, []int64{txn.ScopeID}, mock.Anything,
		).Return(nil).Once()

		// Call the method under test
		err := ModelsService.TransactionModel.UpdateTransaction(context.Background(), txn)
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
		err := ModelsService.TransactionModel.UpdateTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Update Transaction Execution Failure", func(t *testing.T) {
		_, mockM = setupNewMock(t)
		setupForeignKeyMocks(mockM, txn) // Assumes a function to set up foreign key validation

		// Simulate a failure during the execution of the update query
		mockM.ExpectExec("UPDATE transactions").
			WithArgs(txn.SourceID, txn.CategoryID, txn.Amount, txn.Type, txn.Description, txn.ID, txn.ScopeID).
			WillReturnError(sql.ErrConnDone)

		// Call the update method and expect an error
		err := ModelsService.TransactionModel.UpdateTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Handling Transaction Tags Update Fails", func(t *testing.T) {
		_, mockM = setupNewMock(t)
		mockTagModel := new(xmock.MockTagModel)
		mockTransactionTagModel := new(xmock.MockTransactionTagModel)
		ModelsService.TransactionTagModel = mockTransactionTagModel
		ModelsService.TagModel = mockTagModel
		setupForeignKeyMocks(mockM, txn) // Set up foreign key validations

		// Set up mocks for successful transaction update
		mockM.ExpectExec("UPDATE transactions").
			WithArgs(txn.SourceID, txn.CategoryID, txn.Amount, txn.Type, txn.Description, txn.ID, txn.ScopeID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Setup for tag handling
		for _, tag := range txn.Tags {
			mockTagModel.On(
				"GetTagByName",
				mock.Anything, tag, []int64{txn.ScopeID}, mock.Anything,
			).Return(&interfaces.Tag{ID: 1, Name: tag, UserID: txn.UserID, ScopeID: txn.ScopeID}, nil).Once()

			mockTagModel.On(
				"InsertTag",
				mock.Anything, mock.AnythingOfType("*interfaces.Tag"), mock.Anything,
			).Return(nil).Maybe()
		}

		// Simulate a failure in UpdateTagsForTransaction
		mockTransactionTagModel.On(
			"UpdateTagsForTransaction",
			mock.Anything, txn.ID, txn.Tags, []int64{txn.ScopeID}, mock.Anything,
		).Return(sql.ErrConnDone).Once() // Use an appropriate error

		// Call the update method and expect an error
		err := ModelsService.TransactionModel.UpdateTransaction(context.Background(), txn)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
		mockTransactionTagModel.AssertExpectations(t)
	})
}

func TestDeleteTransactionV2(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionModel = NewTransactionModel()

	})
	defer tearDown()

	transactionID := int64(1) // Assume an existing transaction ID for deletion
	scopes := []int64{1}
	db, mockM := setupNewMock(t)
	defer db.Close()

	t.Run("Successful Deletion", func(t *testing.T) {
		// Set up mock for successful deletion
		mockM.ExpectExec("DELETE FROM transactions").
			WithArgs(transactionID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Call the method under test
		err := ModelsService.TransactionModel.DeleteTransaction(context.Background(), transactionID, scopes)
		assert.NoError(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Delete Transaction Execution Failure", func(t *testing.T) {
		_, mockM = setupNewMock(t)

		// Simulate a failure during the execution of the delete query
		mockM.ExpectExec("DELETE FROM transactions").
			WithArgs(transactionID, sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		// Call the delete method and expect an error
		err := ModelsService.TransactionModel.DeleteTransaction(context.Background(), transactionID, scopes)
		assert.Error(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Delete Non-Existent Transaction", func(t *testing.T) {
		_, mockM = setupNewMock(t)

		// Set up mock for a deletion attempt on a non-existent transaction
		mockM.ExpectExec("DELETE FROM transactions").
			WithArgs(transactionID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 0)) // Indicate no rows were affected

		// Call the delete method and check if error or some indication of non-existence is handled
		err := ModelsService.TransactionModel.DeleteTransaction(context.Background(), transactionID, scopes)
		// Assert based on how your application should behave (error or just a no-op)
		// Here it is assumed the method won't error out if no transaction was found
		assert.NoError(t, err)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	// Add more subtests for different edge cases if necessary
	// Example: trying to delete with wrong user ID, etc.
}

func TestGetTransactionByIDV2(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionModel = NewTransactionModel()
		config.UserModel = NewUserModel()
		config.SourceModel = NewSourceModel()
		config.CategoryModel = NewCategoryModel()
		config.TransactionTagModel = NewTransactionTagModel()

	})
	defer tearDown()

	transactionID := int64(1) // Assume some transaction ID
	userID := int64(1)        // Assume some user ID
	scopes := []int64{1}

	// Assume a mock transaction to be returned
	mockTransaction := interfaces.Transaction{
		ID:          transactionID,
		UserID:      userID,
		ScopeID:     1,
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
		rows := sqlmock.NewRows([]string{"transaction_id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description", "scope_id"}).
			AddRow(mockTransaction.ID, mockTransaction.UserID, mockTransaction.SourceID, mockTransaction.CategoryID, mockTransaction.Timestamp, mockTransaction.Amount, mockTransaction.Type, mockTransaction.Description, mockTransaction.ScopeID)

		mockM.ExpectQuery("SELECT (.+) FROM transactions WHERE").WithArgs(transactionID, sqlmock.AnyArg()).WillReturnRows(rows)

		// Assume getTagsForTransaction is properly implemented or mocked
		// Mock the call to getTagsForTransaction if it makes an external call

		transaction, err := ModelsService.TransactionModel.GetTransactionByID(context.Background(), transactionID, scopes)
		assert.NoError(t, err)
		assert.Equal(t, mockTransaction, *transaction)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Transaction Not Found", func(t *testing.T) {
		mockM.ExpectQuery("SELECT (.+) FROM transactions WHERE").WithArgs(transactionID, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows(nil)) // Return empty result set

		transaction, err := ModelsService.TransactionModel.GetTransactionByID(context.Background(), transactionID, scopes)
		assert.Error(t, err) // Assuming function returns error on not found
		assert.Nil(t, transaction)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Query Execution Error", func(t *testing.T) {
		mockM.ExpectQuery("SELECT (.+) FROM transactions WHERE").WithArgs(transactionID, sqlmock.AnyArg()).WillReturnError(sql.ErrConnDone) // Simulate query execution error

		transaction, err := ModelsService.TransactionModel.GetTransactionByID(context.Background(), transactionID, scopes)
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	// Add more subtests if needed, for example to cover the logic within getTagsForTransaction, etc.
}

func TestGetTransactionsByFilterV2(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionModel = NewTransactionModel()
	})
	defer tearDown()

	userID := int64(1) // Assume some user ID
	scopes := []int64{1}
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

		rows := sqlmock.NewRows([]string{"transaction_id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description", "scope_id"})
		for _, txn := range mockTransactions {
			rows = rows.AddRow(txn.ID, txn.UserID, txn.SourceID, txn.CategoryID, txn.Timestamp, txn.Amount, txn.Type, txn.Description, txn.ScopeID)
		}

		mockM.ExpectQuery("SELECT (.+) FROM transactions").WillReturnRows(rows)

		// Assume getTagsForTransaction is properly implemented or mocked

		transactions, err := ModelsService.TransactionModel.GetTransactionsByFilter(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, mockTransactions, transactions)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Query Execution Error", func(t *testing.T) {
		mockM.ExpectQuery("SELECT (.+) FROM transactions").WillReturnError(sql.ErrConnDone) // Simulate query execution error

		transactions, err := ModelsService.TransactionModel.GetTransactionsByFilter(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, transactions)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	t.Run("Row Scan Error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"transaction_id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description", "scope_id"}).
			AddRow(1, userID, 1, 1, time.Now(), 100.0, "expense", "Description", scopes[0])
		mockM.ExpectQuery("SELECT (.+) FROM transactions").WillReturnRows(rows)

		_ = rows.RowError(0, sql.ErrConnDone) // Simulate row scan error on the first row
		transactions, err := ModelsService.TransactionModel.GetTransactionsByFilter(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, transactions)
		assert.NoError(t, mockM.ExpectationsWereMet())
	})

	// Add more subtests if needed, for example to cover different filter criteria and edge cases
}

func TestGetTagsForTransaction(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionModel = NewTransactionModel()
		config.UserModel = NewUserModel()
		config.SourceModel = NewSourceModel()
		config.CategoryModel = NewCategoryModel()
		config.TransactionTagModel = NewTransactionTagModel()
	}) // Set up and initialization
	defer tearDown() // Clean up after the test

	// Assumed transaction details
	transactionID := int64(1)
	userID := int64(1)
	scopes := []int64{1}
	// Mock transaction to use in test
	transaction := &interfaces.Transaction{
		ID:      transactionID,
		UserID:  userID,
		ScopeID: scopes[0],
		// Other necessary fields...
	}

	// Set up the mock database connection and mock expectations
	db, mockM := setupNewMock(t)
	defer db.Close()

	t.Run("Successful Tag Retrieval", func(t *testing.T) {
		// Mock data for tags associated with the transaction
		mockTags := []interfaces.Tag{
			{ID: 1, Name: "Tag1"},
			{ID: 2, Name: "Tag2"},
			// Add more mock tags as necessary
		}

		// Set up SQL mock rows to simulate database response
		rows := sqlmock.NewRows([]string{"transaction_id", "name"})
		for _, tag := range mockTags {
			rows.AddRow(tag.ID, tag.Name)
		}

		// Set up the mock expectation for SQL query
		mockM.ExpectQuery("SELECT t.tag_id, t.name FROM tags t JOIN transaction_tags tt ON t.tag_id = tt.tag_id WHERE tt.transaction_id = ?").
			WithArgs(transactionID).
			WillReturnRows(rows)

		// Call the function under test
		err := getTagsForTransaction(context.Background(), transaction)
		assert.NoError(t, err, "Fetching tags should not produce an error")

		// Check if the tags are correctly assigned to the transaction
		assert.Equal(t, len(mockTags), len(transaction.Tags), "The number of tags should match")
		for i, tag := range mockTags {
			assert.Equal(t, tag.Name, transaction.Tags[i], "Tag names should match")
		}

		// Ensure all expectations were met
		assert.NoError(t, mockM.ExpectationsWereMet(), "All database expectations should be met")
	})

	t.Run("Error Fetching Tags", func(t *testing.T) {
		// Reset mock expectations for a new subtest
		db, mockM = setupNewMock(t)

		// Set up the mock expectation for SQL query with an error
		mockM.ExpectQuery("SELECT t.tag_id, t.name FROM tags t JOIN transaction_tags tt ON t.tag_id = tt.tag_id WHERE tt.transaction_id = ?").
			WithArgs(transactionID).
			WillReturnError(sql.ErrConnDone) // Simulate an error scenario

		// Call the function under test
		err := getTagsForTransaction(context.Background(), transaction)
		assert.Error(t, err, "Fetching tags should produce an error")

		// Ensure all expectations were met
		assert.NoError(t, mockM.ExpectationsWereMet(), "All database expectations should be met")
	})
}

func TestValidateForeignKeyReferences(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionModel = NewTransactionModel()
		config.UserModel = NewUserModel()
		config.SourceModel = NewSourceModel()
		config.CategoryModel = NewCategoryModel()

	})
	defer tearDown()

	txn := interfaces.Transaction{
		UserID:     1,
		SourceID:   1,
		CategoryID: 1,
		ScopeID:    1,
		// ... other necessary fields
	}

	db, mockM := setupNewMock(t)
	defer db.Close()

	t.Run("All References Valid", func(t *testing.T) {
		setupForeignKeyMocks(mockM, txn)

		err := validateForeignKeyReferences(context.Background(), txn)
		assert.NoError(t, err)

		// Ensure all expectations were met (general for sqlmock)
		if err := mockM.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// Additional tests for each failure case: User not found, source not found, category not found
	// Example for user not found:
	t.Run("User Not Found", func(t *testing.T) {
		// Setup new mock database for clean expectation slate
		_, mockM := setupNewMock(t)

		// Mock the user existence check to return no rows, indicating user not found
		mockM.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(txn.UserID).WillReturnRows(sqlmock.NewRows([]string{"1"}))

		err := validateForeignKeyReferences(context.Background(), txn)

		// We are expecting an error because the user is not found
		assert.Error(t, err)
		assert.EqualError(t, err, "user does not exist")

		// Validate that all expectations set on the mock were met
		if err := mockM.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// ... further tests for source not found and category not found ...
}

func TestAddMissingTags(t *testing.T) {
	// Set up the mock environment
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionModel = NewTransactionModel()
	}) // Assume setUp properly initializes mocks and other necessary stuff
	defer tearDown()

	// Assuming these variables are defined and needed for the test
	txn := interfaces.Transaction{
		UserID:     1,
		SourceID:   1,
		CategoryID: 1,
		ScopeID:    1,
		// ... other necessary fields
	}
	ctx := context.Background()
	userID := int64(1)                   // example user ID
	tagNames := []string{"Tag1", "Tag2"} // example tag names
	scopes := []int64{1}                 // example scope IDs
	// Mock TagModel
	mockTagModel := new(xmock.MockTagModel)
	ModelsService.TagModel = mockTagModel

	// Mock the behavior of GetTagByName and InsertTag
	for _, tagName := range tagNames {
		if tagName == "Tag1" {
			// Assume the tag exists
			mockTagModel.On(
				"GetTagByName",
				mock.Anything, // Context
				tagName,
				scopes,
				userID,
				mock.Anything, // Transaction options (if any)
			).Return(&interfaces.Tag{ID: 1, Name: tagName, UserID: userID, ScopeID: scopes[0]}, nil).Once()
		} else {
			// Assume the tag does not exist and needs to be inserted
			mockTagModel.On(
				"GetTagByName",
				mock.Anything,
				tagName,
				scopes,
				userID,
				mock.Anything,
			).Return((*interfaces.Tag)(nil), nil).Once() // Tag not found

			mockTagModel.On(
				"InsertTag",
				mock.Anything,
				mock.MatchedBy(func(tag *interfaces.Tag) bool { return tag.Name == tagName }), // Ensuring the tag name matches
				scopes,
				mock.Anything,
			).Return(nil).Once()
		}
	}

	// Call the function under test
	err := addMissingTags(ctx, txn, nil) // Assuming nil is a valid argument for otx

	// Assertions
	assert.NoError(t, err)

	// Check that all expectations were met
	mockTagModel.AssertExpectations(t)
}
