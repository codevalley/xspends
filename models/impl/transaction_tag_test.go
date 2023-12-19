package impl

import (
	"database/sql"
	"testing"

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
