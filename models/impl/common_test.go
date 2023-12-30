package impl

import (
	"context"
	"os"
	"testing"
	"xspends/models/mock"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/golang/mock/gomock"
)

var (
	ctx          context.Context
	mockExecutor *mock.MockDBExecutor
)

func TestMain(m *testing.M) {
	SQLBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	ctx = context.Background()

	exitVal := m.Run()

	os.Exit(exitVal)
}

func setUp(t *testing.T, modifyConfig func(*ModelsConfig)) func() {
	ctrl := gomock.NewController(t)
	util.InitializeSnowflake()
	// Create mocks for each service
	mockExecutor = mock.NewMockDBExecutor(ctrl)
	mockConfig := &ModelsConfig{
		DBService:           &DBService{Executor: mockExecutor},
		CategoryModel:       new(mock.MockCategoryModel),
		SourceModel:         new(mock.MockSourceModel),
		UserModel:           new(mock.MockUserModel),
		TagModel:            new(mock.MockTagModel),
		TransactionTagModel: new(mock.MockTransactionTagModel),
		TransactionModel:    new(mock.MockTransactionModel),
	}

	// Allow tests to modify the mock configuration as needed
	if modifyConfig != nil {
		modifyConfig(mockConfig)
	}
	isTesting = true
	// Initialize ModelsService with the (potentially modified) mock configuration
	InitModelsService(mockConfig)

	// Return a function to teardown the mock controller after the test
	return func() {
		ctrl.Finish() // Ensure all mock expectations are met
		// Any additional cleanup can be added here
	}
}
