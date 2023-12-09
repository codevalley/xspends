package testutils

import (
	"context"
	"testing"

	"xspends/models/impl"
	"xspends/models/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

func SetupModelTestEnvironment(t *testing.T) (context.Context, *impl.ModelsServiceContainer, *mock.MockDBExecutor, sqlmock.Sqlmock, *mock.MockCategoryModel, func()) {
	ctx := context.Background()

	// Set up gomock controller and mock executor
	ctrl := gomock.NewController(t)
	mockExecutor := mock.NewMockDBExecutor(ctrl)

	// Set up sqlmock
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %s", err)
	}

	// Create mock DBService
	mockDBService := &impl.DBService{Executor: mockExecutor}

	// Create a mock CategoryModel
	mockCategoryModel := new(mock.MockCategoryModel)

	// Create ModelsServiceContainer with mocks
	mockModelService := &impl.ModelsServiceContainer{
		DBService:     mockDBService,
		CategoryModel: mockCategoryModel,
		// Initialize other services as necessary
	}

	// Teardown function
	tearDown := func() {
		ctrl.Finish()
		db.Close()
	}

	return ctx, mockModelService, mockExecutor, sqlMock, mockCategoryModel, tearDown
}
