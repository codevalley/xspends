package impl

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInsertTransactionTag(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	transactionID := int64(1)
	tagID := int64(1)

	mockExecutor.EXPECT().
		ExecContext(ctx, gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := mockModelService.TransactionTagModel.InsertTransactionTag(ctx, transactionID, tagID)
	assert.NoError(t, err)

	//test for generic query error
	mockExecutor.EXPECT().
		ExecContext(ctx, gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), errors.New("failed to build SQL query for InsertTransactionTag")).
		Times(1)

	err = mockModelService.TransactionTagModel.InsertTransactionTag(ctx, transactionID, tagID)
	assert.Error(t, err)

}

func TestDeleteTransactionTag(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	transactionID := int64(1)
	tagID := int64(1)

	mockExecutor.EXPECT().
		ExecContext(ctx, gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := mockModelService.TransactionTagModel.DeleteTransactionTag(ctx, transactionID, tagID)
	assert.NoError(t, err)

	//test for generic query error
	mockExecutor.EXPECT().
		ExecContext(ctx, gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), errors.New("error deleting transaction tag")).
		Times(1)

	err = mockModelService.TransactionTagModel.DeleteTransactionTag(ctx, transactionID, tagID)
	assert.Error(t, err)
}

func TestGetTagsByTransactionID(t *testing.T) {
	tearDown := setUp(t)
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
	mockModelService.DBService.Executor = db

	transactionID := int64(1)

	// Set up expectations for QueryContext to return mock rows
	mockRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Tag1").
		AddRow(2, "Tag2")

	// Update the expected SQL query pattern to match the actual query
	expectedSQLPattern := `SELECT t\.id, t\.name FROM tags t JOIN transaction_tags tt ON t\.id = tt\.tag_id WHERE tt\.transaction_id = \?`

	mock.ExpectQuery(expectedSQLPattern).WithArgs(transactionID).
		WillReturnRows(mockRows)

	tags, err := mockModelService.TransactionTagModel.GetTagsByTransactionID(ctx, transactionID)
	assert.NoError(t, err)
	assert.NotNil(t, tags)

	// Add assertions to validate the returned tags

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestUpdateTagsForTransaction(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	transactionID := int64(1)
	userID := int64(1)
	tags := []string{"Tag1", "Tag2"}

	// Create a sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred when creating a mock database connection: %v", err)
	}
	defer db.Close()

	// Replace the DBService Executor with the mock db
	mockModelService.DBService.Executor = db

	// Mocking the DeleteTagsFromTransaction SQL query
	mock.ExpectExec("DELETE FROM transaction_tags WHERE transaction_id = ?").
		WithArgs(transactionID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // assuming 1 row affected

	// Mocking the GetTagByName and InsertTransactionTag SQL query for each tag
	for _, tagName := range tags {
		tagID := int64(1) // Mock tag ID

		// Adjust the expected SQL query pattern to match the actual query
		expectedSQLPattern := `SELECT id, user_id, name, created_at, updated_at FROM tags WHERE name = \? AND user_id = \?`

		mock.ExpectQuery(expectedSQLPattern).
			WithArgs(tagName, userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).AddRow(tagID, userID, tagName, time.Now(), time.Now()))

		// Mocking InsertTransactionTag
		mock.ExpectExec("INSERT INTO transaction_tags").
			WithArgs(transactionID, tagID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1)) // assuming 1 row affected
	}

	err = mockModelService.TransactionTagModel.UpdateTagsForTransaction(ctx, transactionID, tags, userID)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestAddTagsToTransaction(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	transactionID := int64(1)
	userID := int64(1)
	tags := []string{"Tag1", "Tag2"}

	// Create a sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred when creating a mock database connection: %v", err)
	}
	defer db.Close()

	// Replace the DBService Executor with the mock db
	mockModelService.DBService.Executor = db

	// Mocking the GetTagByName and InsertTransactionTag SQL query for each tag
	for _, tagName := range tags {
		tagID := int64(1) // Mock tag ID

		// Mocking GetTagByName
		mock.ExpectQuery(`SELECT id, user_id, name, created_at, updated_at FROM tags WHERE name = \? AND user_id = \?`).
			WithArgs(tagName, userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).AddRow(tagID, userID, tagName, time.Now(), time.Now()))

		// Mocking InsertTransactionTag
		mock.ExpectExec("INSERT INTO transaction_tags").
			WithArgs(transactionID, tagID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1)) // assuming 1 row affected
	}

	err = mockModelService.TransactionTagModel.AddTagsToTransaction(ctx, transactionID, tags, userID)
	assert.NoError(t, err)

	// Testing error scenario when GetTagByName fails
	mock.ExpectQuery(`SELECT id, user_id, name, created_at, updated_at FROM tags WHERE name = \? AND user_id = \?`).
		WithArgs(tags[0], userID).
		WillReturnError(errors.New("error getting tag by name"))

	err = mockModelService.TransactionTagModel.AddTagsToTransaction(ctx, transactionID, tags[:1], userID)
	assert.Error(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestDeleteTagsFromTransaction(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	transactionID := int64(1)

	// Create a sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred when creating a mock database connection: %v", err)
	}
	defer db.Close()

	// Replace the DBService Executor with the mock db
	mockModelService.DBService.Executor = db

	// Mocking the DeleteTagsFromTransaction SQL query
	mock.ExpectExec("DELETE FROM transaction_tags WHERE transaction_id = ?").
		WithArgs(transactionID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // assuming 1 row affected

	err = mockModelService.TransactionTagModel.DeleteTagsFromTransaction(ctx, transactionID)
	assert.NoError(t, err)

	// Testing error scenario for the SQL query
	mock.ExpectExec("DELETE FROM transaction_tags WHERE transaction_id = ?").
		WithArgs(transactionID).
		WillReturnError(errors.New("error deleting tags from transaction"))

	err = mockModelService.TransactionTagModel.DeleteTagsFromTransaction(ctx, transactionID)
	assert.Error(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}
