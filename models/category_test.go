package models

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"
	"xspends/models/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockDBService *DBService
	mockExecutor  *mock.MockDBExecutor
	ctx           context.Context
)

func TestMain(m *testing.M) {
	// Set up that is common to all tests
	SQLBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	ctx = context.Background()

	// Run the tests
	exitVal := m.Run()

	// Teardown if necessary

	os.Exit(exitVal)
}

func setUp(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockExecutor = mock.NewMockDBExecutor(ctrl)
	mockDBService = &DBService{Executor: mockExecutor}
	return func() { ctrl.Finish() }
}

func TestInsertCategory(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	category := &Category{
		UserID:      1,
		Name:        "Test Category",
		Description: "This is a test category",
		Icon:        "test-icon",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sqlmock.NewResult(1, 1), nil).
		Times(1)

	err := InsertCategory(ctx, category, mockDBService)
	assert.NoError(t, err)
}

func TestUpdateCategory(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	category := &Category{
		ID:          1,
		UserID:      1,
		Name:        "Updated Category",
		Description: "Updated Description",
		Icon:        "updated-icon",
		UpdatedAt:   time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := UpdateCategory(ctx, category, mockDBService)
	assert.NoError(t, err)
}

func TestDeleteCategory(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	categoryID := int64(1)
	userID := int64(1)

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := DeleteCategory(ctx, categoryID, userID, mockDBService)
	assert.NoError(t, err)
}

func TestGetAllCategories(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	userID := int64(1)

	// Here we use sqlmock to simulate database response
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "icon", "created_at", "updated_at"}).
		AddRow(1, userID, "Category 1", "Description 1", "icon1", time.Now(), time.Now()).
		AddRow(2, userID, "Category 2", "Description 2", "icon2", time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnRows(rows)

	categories, err := GetAllCategories(ctx, userID, mockDBService)
	assert.NoError(t, err)
	assert.Len(t, categories, 2)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// ... Additional tests for GetCategoryByID, GetPagedCategories, and CategoryIDExists ...
