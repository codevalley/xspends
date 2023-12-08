package models

import (
	"context"
	"os"
	"testing"
	"time"
	"xspends/models/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
	ctrl := gomock.NewController(t)
	mockExecutor = mock.NewMockDBExecutor(ctrl)
	mockDBService = &DBService{Executor: mockExecutor}
	mockModelService = &ModelsServiceContainer{
		DBService:     mockDBService,
		CategoryModel: &CategoryModel{}}
	ModelsService = mockModelService
	return func() { ctrl.Finish() }
}

// Insert your test functions here, following the pattern demonstrated below.

func TestInsertCategoryWithInvalidInput(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	cm := &CategoryModel{}

	invalidCategories := []*Category{
		{UserID: 0, Name: "Test Category", Description: "Description"},
		{UserID: 1, Name: "", Description: "Description"},
		{UserID: 1, Name: "Test Category", Description: string(make([]byte, 501))},
	}

	for _, invalidCategory := range invalidCategories {
		err := cm.InsertCategory(ctx, invalidCategory, nil)
		assert.EqualError(t, err, ErrInvalidInput)
	}
}

// TestInsertCategoryWithDatabaseError verifies that the function returns an error for database errors.
func TestInsertCategoryWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	cm := &CategoryModel{}
	category := &Category{UserID: 1, Name: "Test Category", Description: "Description"}
	err := cm.InsertCategory(ctx, category, nil)
	assert.EqualError(t, err, "executing insert statement failed: database error")
}

// TestUpdateCategoryWithDatabaseError verifies that the function returns an error for database errors.
func TestUpdateCategoryWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	cm := &CategoryModel{}
	category := &Category{ID: 1, UserID: 1, Name: "Updated Category", Description: "Updated Description"}
	err := cm.UpdateCategory(ctx, category, nil)
	assert.EqualError(t, err, "executing update statement failed: database error")
}

// TestDeleteCategoryWithDatabaseError verifies that the function returns an error for database errors.
func TestDeleteCategoryWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	cm := &CategoryModel{}
	err := cm.DeleteCategory(ctx, 1, 1, nil)
	assert.EqualError(t, err, "executing delete statement failed: database error")
}

// TestGetAllCategoriesWithDatabaseError verifies that the function returns an error for database errors.
func TestGetAllCategoriesWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnError(errors.New("database error"))

	cm := &CategoryModel{}
	_, err = cm.GetAllCategories(ctx, 1, nil)
	assert.EqualError(t, err, "querying categories failed: database error")
}

// TestGetCategoryByIDWithCategoryNotFound verifies that the function returns an error for a non-existent category.
func TestGetCategoryByIDWithCategoryNotFound(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnRows(sqlmock.NewRows([]string{}))

	cm := &CategoryModel{}
	_, err = cm.GetCategoryByID(ctx, 1, 1, nil)
	assert.EqualError(t, err, "category not found")
}

// TestGetCategoryByIDWithDatabaseError verifies that the function returns an error for database errors.
func TestGetCategoryByIDWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnError(errors.New("database error"))

	cm := &CategoryModel{}
	_, err = cm.GetCategoryByID(ctx, 1, 1, nil)
	assert.EqualError(t, err, "querying category by ID failed: database error")
}

// TestGetPagedCategoriesWithDatabaseError verifies that the function returns an error for database errors.
func TestGetPagedCategoriesWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnError(errors.New("database error"))

	cm := &CategoryModel{}
	_, err = cm.GetPagedCategories(ctx, 1, 10, 1, nil)
	assert.EqualError(t, err, "querying paginated categories failed: database error")
}

// TestCategoryIDExistsWithDatabaseError verifies that the function returns an error for database errors.
func TestCategoryIDExistsWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnError(errors.New("database error"))

	cm := &CategoryModel{}
	exists, err := cm.CategoryIDExists(ctx, 1, 1, nil)
	assert.EqualError(t, err, "checking category existence failed: database error")
	assert.False(t, exists)
}

// TestGetCategoryByIDWithEmptyIcon verifies that the category is returned with an empty icon when not set.
func TestGetCategoryByIDWithEmptyIcon(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	expectedCategory := &Category{
		ID:          1,
		UserID:      1,
		Name:        "Test Category",
		Description: "Description",
		Icon:        "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "icon", "created_at", "updated_at"}).
		AddRow(expectedCategory.ID, expectedCategory.UserID, expectedCategory.Name, expectedCategory.Description, "", expectedCategory.CreatedAt, expectedCategory.UpdatedAt)

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnRows(rows)

	cm := &CategoryModel{}
	category, err := cm.GetCategoryByID(ctx, 1, 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, category)
}

// TestGetAllCategoriesWithEmptyResults verifies that an empty slice is returned when there are no categories.
func TestGetAllCategoriesWithEmptyResults(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	rows := sqlmock.NewRows([]string{})

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnRows(rows)

	cm := &CategoryModel{}
	categories, err := cm.GetAllCategories(ctx, 1, nil)
	assert.NoError(t, err)
	assert.Empty(t, categories)
}

// TestCategoryIDExistsWithNonExistentCategory verifies that false is returned for non-existent category.
func TestCategoryIDExistsWithNonExistentCategory(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	categoryID := int64(1)
	userID := int64(1)

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	cm := &CategoryModel{}
	actualExists, err := cm.CategoryIDExists(ctx, categoryID, userID, nil)
	assert.NoError(t, err)
	assert.False(t, actualExists)
}

// TestDeleteCategoryWithDatabase verifies that DeleteCategory works with a mock database.
func TestDeleteCategoryWithDatabase(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	categoryID := int64(1)
	userID := int64(1)

	mock.ExpectExec("^DELETE FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	cm := &CategoryModel{}
	err = cm.DeleteCategory(ctx, categoryID, userID, nil)
	assert.NoError(t, err)
}

// TestGetCategoryByIDWithDatabase tests retrieval of a category by ID using a mock database.
func TestGetCategoryByIDWithDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}
	mockModelService = &ModelsServiceContainer{
		DBService:     mockDBService,
		CategoryModel: &CategoryModel{}}
	ModelsService = mockModelService

	categoryID := int64(1)
	userID := int64(1)

	expectedCategory := &Category{
		ID:          categoryID,
		UserID:      userID,
		Name:        "Test Category",
		Description: "Description",
		Icon:        "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "icon", "created_at", "updated_at"}).
		AddRow(expectedCategory.ID, expectedCategory.UserID, expectedCategory.Name, expectedCategory.Description, expectedCategory.Icon, expectedCategory.CreatedAt, expectedCategory.UpdatedAt)

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnRows(rows)

	cm := &CategoryModel{}
	category, err := cm.GetCategoryByID(ctx, categoryID, userID, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, category)
}

// TestCategoryIDExistsWithDatabase tests checking of category existence using a mock database.
func TestCategoryIDExistsWithDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}
	mockModelService = &ModelsServiceContainer{
		DBService:     mockDBService,
		CategoryModel: &CategoryModel{}}
	ModelsService = mockModelService

	categoryID := int64(1)
	userID := int64(1)

	exists := true

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(categoryID))

	cm := &CategoryModel{}
	actualExists, err := cm.CategoryIDExists(ctx, categoryID, userID, nil)
	assert.NoError(t, err)
	assert.Equal(t, exists, actualExists)
}

// TestGetPagedCategoriesWithEmptyResults verifies that an empty slice is returned when there are no categories for pagination.
func TestGetPagedCategoriesWithEmptyResults(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	rows := sqlmock.NewRows([]string{})

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnRows(rows)

	cm := &CategoryModel{}
	categories, err := cm.GetPagedCategories(ctx, 1, 10, 1, nil)
	assert.NoError(t, err)
	assert.Empty(t, categories)
}

// TestUpdateCategoryWithValidInput verifies that UpdateCategory successfully updates a category.
func TestUpdateCategoryWithValidInput(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	category := &Category{ID: 1, UserID: 1, Name: "Updated Category", Description: "Updated Description"}

	mock.ExpectExec("^UPDATE categories SET").WithArgs(category.Name, category.Description, category.Icon, sqlmock.AnyArg(), category.ID, category.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

	cm := &CategoryModel{}
	err = cm.UpdateCategory(ctx, category, nil)
	assert.NoError(t, err)
}

// Add any additional test cases as required for your application...
