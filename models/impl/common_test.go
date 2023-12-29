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
	mockModelService = &ModelsServiceContainer{
		DBService:           mockDBService,
		CategoryModel:       &CategoryModel{},
		SourceModel:         &SourceModel{},
		UserModel:           &UserModel{},
		TagModel:            &TagModel{},
		TransactionTagModel: &TransactionTagModel{},
		TransactionModel:    &TransactionModel{},
	}
	dbService = mockDBService
	ModelsService = mockModelService

	return func() { ctrl.Finish() }
}
