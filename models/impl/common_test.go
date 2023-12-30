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
	ctx              context.Context
	mockExecutor     *mock.MockDBExecutor
	mockDBService    *DBService
	mockModelService *ModelsServiceContainer
)

func TestMain(m *testing.M) {
	SQLBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	ctx = context.Background()

	exitVal := m.Run()

	os.Exit(exitVal)
}

func setUp(t *testing.T) func() {
	util.InitializeSnowflake()
	ctrl := gomock.NewController(t)
	mockExecutor = mock.NewMockDBExecutor(ctrl)
	mockDBService = &DBService{Executor: mockExecutor} //bad code, this variable is in db.go

	mockConfig := &ModelsConfig{
		DBService:           mockDBService,
		CategoryModel:       new(mock.MockCategoryModel),
		SourceModel:         new(mock.MockSourceModel),
		UserModel:           new(mock.MockUserModel),
		TagModel:            new(mock.MockTagModel),
		TransactionTagModel: new(mock.MockTransactionTagModel),
		TransactionModel:    new(mock.MockTransactionModel),
	}
	// Initialize the models service with the mock configuration
	InitModelsService(mockConfig)

	return func() { ctrl.Finish() }
}
