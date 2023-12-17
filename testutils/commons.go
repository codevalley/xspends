package testutils

import (
	"context"
	"testing"

	"xspends/models/impl"
	"xspends/models/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

var (
	MockDBService           *impl.DBService
	MockCategoryModel       *mock.MockCategoryModel
	MockSourceModel         *mock.MockSourceModel
	MockUserModel           *mock.MockUserModel
	MockTagModel            *mock.MockTagModel
	MockTransactionTagModel *mock.MockTransactionTagModel
)

func SetupModelTestEnvironment(t *testing.T) (context.Context, *impl.ModelsServiceContainer, *mock.MockDBExecutor, sqlmock.Sqlmock, func()) {
	ctx := context.TODO()

	// Set up gomock controller and mock executor
	ctrl := gomock.NewController(t)
	mockExecutor := mock.NewMockDBExecutor(ctrl)

	// Set up sqlmock
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %s", err)
	}

	// Create mock DBService
	MockDBService = &impl.DBService{Executor: mockExecutor}

	// Create a mock CategoryModel
	MockCategoryModel = new(mock.MockCategoryModel)
	MockUserModel = new(mock.MockUserModel)
	MockSourceModel = new(mock.MockSourceModel)
	MockTagModel = new(mock.MockTagModel)
	// Create ModelsServiceContainer with mocks
	impl.ModelsService = &impl.ModelsServiceContainer{
		DBService:           MockDBService,
		CategoryModel:       MockCategoryModel,
		UserModel:           MockUserModel,
		SourceModel:         MockSourceModel,
		TagModel:            MockTagModel,
		TransactionTagModel: MockTransactionTagModel,
		// Initialize other services as necessary
	}

	// Teardown function
	tearDown := func() {
		ctrl.Finish()
		db.Close()
	}

	return ctx, impl.ModelsService, mockExecutor, sqlMock, tearDown
}
