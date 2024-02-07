package impl

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInsertTransactionTag(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionTagModel = NewTransactionTagModel()
	})
	defer tearDown()

	transactionID := int64(1)
	tagID := int64(1)

	mockExecutor.EXPECT().
		ExecContext(ctx, gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := ModelsService.TransactionTagModel.InsertTransactionTag(ctx, transactionID, tagID)
	assert.NoError(t, err)

	//test for generic query error
	mockExecutor.EXPECT().
		ExecContext(ctx, gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), errors.New("failed to build SQL query for InsertTransactionTag")).
		Times(1)

	err = ModelsService.TransactionTagModel.InsertTransactionTag(ctx, transactionID, tagID)
	assert.Error(t, err)

}

func TestDeleteTransactionTag(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionTagModel = NewTransactionTagModel()
	})
	defer tearDown()

	transactionID := int64(1)
	tagID := int64(1)

	mockExecutor.EXPECT().
		ExecContext(ctx, gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := ModelsService.TransactionTagModel.DeleteTransactionTag(ctx, transactionID, tagID)
	assert.NoError(t, err)

	//test for generic query error
	mockExecutor.EXPECT().
		ExecContext(ctx, gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), errors.New("error deleting transaction tag")).
		Times(1)

	err = ModelsService.TransactionTagModel.DeleteTransactionTag(ctx, transactionID, tagID)
	assert.Error(t, err)
}

func TestGetTagsByTransactionID(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionTagModel = NewTransactionTagModel()
	})
	defer tearDown()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred when creating a mock database connection: %v", err)
	}
	defer db.Close()

	// Replace your mockExecutor with db
	ModelsService.DBService.Executor = db

	transactionID := int64(1)

	// Set up expectations for QueryContext to return mock rows
	mockRows := sqlmock.NewRows([]string{"tag_id", "name"}).
		AddRow(1, "Tag1").
		AddRow(2, "Tag2")

	//TODO: Remove hardcoded literals in SQL queries
	// Update the expected SQL query pattern to match the actual query
	expectedSQLPattern := `SELECT tag_id, name FROM tags t JOIN transaction_tags tt ON t\.tag_id = tt\.tag_id WHERE tt\.transaction_id = \?`

	mock.ExpectQuery(expectedSQLPattern).WithArgs(transactionID).
		WillReturnRows(mockRows)

	tags, err := ModelsService.TransactionTagModel.GetTagsByTransactionID(ctx, transactionID)
	assert.NoError(t, err)
	assert.NotNil(t, tags)

	// Add assertions to validate the returned tags

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestUpdateTagsForTransaction(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionTagModel = NewTransactionTagModel()
		config.TagModel = NewTagModel()
	})
	defer tearDown()

	// Create a sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred when creating a mock database connection: %v", err)
	}
	defer db.Close()

	// Replace the DBService Executor with the mock db
	ModelsService.DBService.Executor = db
	transactionID := int64(1)
	userID := int64(1)
	scopes := []int64{1, 2, 3} // Assuming multiple scopes for demonstration
	tags := []string{"Tag1", "Tag2"}

	// Mocking the DeleteTagsFromTransaction SQL query
	mock.ExpectExec("DELETE FROM transaction_tags WHERE transaction_id = ?").
		WithArgs(transactionID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // assuming 1 row affected

	// Mocking the GetTagByName and InsertTransactionTag SQL query for each tag
	for _, tagName := range tags {
		tagID := int64(1) // Mock tag ID for simplicity

		// Correctly handling IN clause with multiple values
		// Note: You might need to adjust this based on how your actual application constructs the query
		mock.ExpectQuery(`SELECT tag_id, user_id, name, scope_id, created_at, updated_at FROM tags WHERE name = \? AND scope_id IN \(\?,\?,\?\)`).
			WithArgs(tagName, scopes[0], scopes[1], scopes[2]).
			WillReturnRows(sqlmock.NewRows([]string{"tag_id", "user_id", "name", "scope_id", "created_at", "updated_at"}).
				AddRow(tagID, userID, tagName, scopes[0], time.Now(), time.Now()))

		// Mocking InsertTransactionTag for each tag
		// This part of your test seems correct; adjust if necessary based on actual logic
		mock.ExpectExec("INSERT INTO transaction_tags").
			WithArgs(transactionID, tagID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1)) // assuming 1 row affected
	}

	// Execute the method under test
	err = ModelsService.TransactionTagModel.UpdateTagsForTransaction(context.Background(), transactionID, tags, scopes)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestAddTagsToTransaction(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionTagModel = NewTransactionTagModel()
		config.TagModel = NewTagModel()
	})
	defer tearDown()
	// Create a sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred when creating a mock database connection: %v", err)
	}
	defer db.Close()

	// Replace the DBService Executor with the mock db
	ModelsService.DBService.Executor = db

	transactionID := int64(1)
	userID := int64(1)
	scopes := []int64{1}
	tags := []string{"Tag1", "Tag2"}

	// Mocking the GetTagByName and InsertTransactionTag SQL query for each tag
	for _, tagName := range tags {
		tagID := int64(1) // Mock tag ID

		// Mocking GetTagByName
		mock.ExpectQuery(`SELECT tag_id, user_id, name, scope_id, created_at, updated_at FROM tags WHERE name = \? AND scope_id IN (\?)`).
			WithArgs(tagName, scopes[0]).
			WillReturnRows(sqlmock.NewRows([]string{"tag_id", "user_id", "name", "scope_id", "created_at", "updated_at"}).
				AddRow(tagID, userID, tagName, scopes[0], time.Now(), time.Now()))

		// Mocking InsertTransactionTag
		mock.ExpectExec("INSERT INTO transaction_tags").
			WithArgs(transactionID, tagID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1)) // assuming 1 row affected
	}

	err = ModelsService.TransactionTagModel.AddTagsToTransaction(ctx, transactionID, tags, scopes)
	assert.NoError(t, err)

	// Testing error scenario when GetTagByName fails
	mock.ExpectQuery(`SELECT tag_id, user_id, name, created_at, updated_at FROM tags WHERE name = \? AND scope_id IN (\?)`).
		WithArgs(tags[0], scopes[0]).
		WillReturnError(errors.New("error getting tag by name"))

	err = ModelsService.TransactionTagModel.AddTagsToTransaction(ctx, transactionID, tags[:1], scopes)
	assert.Error(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestDeleteTagsFromTransaction(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TransactionTagModel = NewTransactionTagModel()
		config.TagModel = NewTagModel()
	})
	defer tearDown()

	transactionID := int64(1)

	// Create a sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred when creating a mock database connection: %v", err)
	}
	defer db.Close()

	// Replace the DBService Executor with the mock db
	ModelsService.DBService.Executor = db

	// Mocking the DeleteTagsFromTransaction SQL query
	mock.ExpectExec("DELETE FROM transaction_tags WHERE transaction_id = ?").
		WithArgs(transactionID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // assuming 1 row affected

	err = ModelsService.TransactionTagModel.DeleteTagsFromTransaction(ctx, transactionID)
	assert.NoError(t, err)

	// Testing error scenario for the SQL query
	mock.ExpectExec("DELETE FROM transaction_tags WHERE transaction_id = ?").
		WithArgs(transactionID).
		WillReturnError(errors.New("error deleting tags from transaction"))

	err = ModelsService.TransactionTagModel.DeleteTagsFromTransaction(ctx, transactionID)
	assert.Error(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}
