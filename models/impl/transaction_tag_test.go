package impl

import (
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInsertTransactionTag(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	transactionID := int64(1)
	tagID := int64(1)

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := mockModelService.TransactionTagModel.InsertTransactionTag(ctx, transactionID, tagID)
	assert.NoError(t, err)

}
