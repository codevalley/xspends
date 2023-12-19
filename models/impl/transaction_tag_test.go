package impl

import (
	"database/sql"
	"testing"

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
