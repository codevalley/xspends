package impl

import (
	"testing"
	"xspends/models/interfaces"
	xmock "xspends/models/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsertTransaction(t *testing.T) {
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

	// Create a sqlmock database connection
	db, mockM, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred when creating a mock database connection: %v", err)
	}
	defer db.Close()

	// Replace the DBService Executor with the mock db
	mockModelService.DBService.Executor = db
	mockModelService.TransactionTagModel = new(xmock.MockTransactionTagModel)
	mockModelService.TagModel = new(xmock.MockTagModel)
	// Correctly return rows for foreign key validation queries
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

	// Expectation for the main INSERT operation into transactions
	mockM.ExpectExec("INSERT INTO transactions").
		WithArgs(sqlmock.AnyArg(), txn.UserID, txn.SourceID, txn.CategoryID, sqlmock.AnyArg(), txn.Amount, txn.Type, txn.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tagModelMock := mockModelService.TagModel.(*xmock.MockTagModel)
	transactionTagModelMock := mockModelService.TransactionTagModel.(*xmock.MockTransactionTagModel)

	// Setting expectation for GetTagByName for each tag
	for _, tag := range txn.Tags {
		// If the tag is found, return a tag object; if not, return nil.
		// Adjust the second return value based on the actual method signature.
		// Here, it's assumed that the method returns (*Tag, error).
		tagModelMock.On(
			"GetTagByName", // method name
			mock.Anything,  // match any context
			tag,            // the specific tag name
			txn.UserID,     // the specific user ID
			mock.Anything,  // match any transaction options
		).Return(&interfaces.Tag{ID: 1, UserID: txn.UserID, Name: tag}, nil) // Adjust return values as needed
	}

	// Setting expectation for AddTagsToTransaction for the transaction
	transactionTagModelMock.On(
		"AddTagsToTransaction", // method name
		mock.Anything,          // match any context
		mock.Anything,          // match any transaction ID
		mock.Anything,          // match any tags
		txn.UserID,             // the specific user ID
		mock.Anything,          // match any transaction options
	).Return(nil) // Adjust based on actual method signature

	// Call the method under test
	err = mockModelService.TransactionModel.InsertTransaction(ctx, txn)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mockM.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}
