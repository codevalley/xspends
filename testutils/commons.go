package testutils

import (
	"context"
	"testing"

	"xspends/models/impl"
	"xspends/models/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
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

	// Create a mock CategoryModel
	mockCategoryModel := new(mock.MockCategoryModel)
	mockUserModel := new(mock.MockUserModel)
	mockSourceModel := new(mock.MockSourceModel)
	mockTagModel := new(mock.MockTagModel)
	mockTransactionTagModel := new(mock.MockTransactionTagModel)
	mockTransactionModel := new(mock.MockTransactionModel)
	mockScopeModel := new(mock.MockScopeModel)
	mockGroupModel := new(mock.MockGroupModel)
	//create mockconfigs
	mockConfig := &impl.ModelsConfig{
		DBService:           &impl.DBService{Executor: mockExecutor},
		CategoryModel:       mockCategoryModel,
		SourceModel:         mockSourceModel,
		UserModel:           mockUserModel,
		TagModel:            mockTagModel,
		TransactionTagModel: mockTransactionTagModel,
		TransactionModel:    mockTransactionModel,
		ScopeModel:          mockScopeModel,
		GroupModel:          mockGroupModel,
	}
	// Initialize ModelsService with mock configuration
	impl.InitModelsService(mockConfig)
	// Teardown function
	tearDown := func() {
		ctrl.Finish()
		db.Close()
	}

	return ctx, impl.ModelsService, mockExecutor, sqlMock, tearDown
}
