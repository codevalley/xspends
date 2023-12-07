package models

import (
	"context"
	"errors"
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

// TestInsertCategoryWithInvalidInput verifies that the function returns an error for invalid input.
func TestInsertCategoryWithInvalidInput(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	invalidCategories := []*Category{
		{UserID: 0, Name: "Test Category", Description: "Description"},             // User ID must be positive.
		{UserID: 1, Name: "", Description: "Description"},                          // Name cannot be empty.
		{UserID: 1, Name: "Test Category", Description: string(make([]byte, 501))}, // Description exceeds max length.
	}

	for _, invalidCategory := range invalidCategories {
		err := InsertCategory(ctx, invalidCategory, mockDBService)
		assert.EqualError(t, err, ErrInvalidInput)
	}
}

// TestInsertCategoryWithDatabaseError verifies that the function returns an error for database errors.
func TestInsertCategoryWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	category := &Category{UserID: 1, Name: "Test Category", Description: "Description"}
	err := InsertCategory(ctx, category, mockDBService)
	assert.EqualError(t, err, "executing insert statement failed: database error")
}

// TestUpdateCategoryWithDatabaseError verifies that the function returns an error for database errors.
func TestUpdateCategoryWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	category := &Category{ID: 1, UserID: 1, Name: "Updated Category", Description: "Updated Description"}
	err := UpdateCategory(ctx, category, mockDBService)
	assert.EqualError(t, err, "executing update statement failed: database error")
}

// TestDeleteCategoryWithDatabaseError verifies that the function returns an error for database errors.
func TestDeleteCategoryWithDatabaseError(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	err := DeleteCategory(ctx, 1, 1, mockDBService)
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

	_, err = GetAllCategories(ctx, 1, mockDBService)
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

	_, err = GetCategoryByID(ctx, 1, 1, mockDBService)
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

	_, err = GetCategoryByID(ctx, 1, 1, mockDBService)
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

	_, err = GetPagedCategories(ctx, 1, 10, 1, mockDBService)
	assert.EqualError(t, err, "querying paginated categories failed: database error")
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

	category, err := GetCategoryByID(ctx, 1, 1, mockDBService)
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

	categories, err := GetAllCategories(ctx, 1, mockDBService)
	assert.NoError(t, err)
	assert.Empty(t, categories)
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

	exists, err := CategoryIDExists(ctx, 1, 1, mockDBService)
	assert.EqualError(t, err, "checking category existence failed: database error")
	assert.False(t, exists)
}

func TestGetCategoryByIDWithDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}

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

	category, err := GetCategoryByID(ctx, categoryID, userID, mockDBService)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, category)
}

func TestCategoryIDExistsWithNonExistentCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}

	categoryID := int64(1)
	userID := int64(1)

	exists := false

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	actualExists, err := CategoryIDExists(ctx, categoryID, userID, mockDBService)
	assert.NoError(t, err)
	assert.Equal(t, exists, actualExists)
}
func TestCategoryIDExistsWithDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}

	categoryID := int64(1)
	userID := int64(1)

	exists := true

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(categoryID))

	actualExists, err := CategoryIDExists(ctx, categoryID, userID, mockDBService)
	assert.NoError(t, err)
	assert.Equal(t, exists, actualExists)
}

func TestDeleteCategoryWithDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}

	categoryID := 1
	userID := 1

	mock.ExpectExec("^DELETE FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = DeleteCategory(ctx, int64(categoryID), int64(userID), mockDBService)
	assert.NoError(t, err)
	assert.NoError(t, err)
}

/*


// TestGetPagedCategoriesWithInvalidPage verifies that the function returns an error for an invalid page number.
func TestGetPagedCategoriesWithInvalidPage(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	_, err := GetPagedCategories(ctx, 0, 10, 1, mockDBService)
	assert.EqualError(t, err, "invalid page number: page number must be positive")
}

// TestUpdateCategoryWithCategoryNotFound verifies that the function returns an error for a non-existent category.
func TestUpdateCategoryWithCategoryNotFound(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	// Expect ExecContext to be called once and return no rows affected
	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(sql.Result(nil), nil).Times(1)

	category := &Category{ID: 1, UserID: 1, Name: "Updated Category", Description: "Updated Description"}

	// UpdateCategory should return an error indicating category not found
	err := UpdateCategory(ctx, category, mockDBService)
	assert.Error(t, err)

	assert.EqualError(t, err, "category not found")
}

// TestGetPagedCategoriesWithLastPage verifies that the correct page is returned when fetching the last page.
func TestGetPagedCategoriesWithLastPage(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService.Executor = db

	// totalCategories := 3
	itemsPerPage := 2
	expectedCategories := []*Category{
		{ID: 2, UserID: 1, Name: "Category 2", Description: "Description 2", Icon: "icon2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, UserID: 1, Name: "Category 3", Description: "Description 3", Icon: "icon3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "icon", "created_at", "updated_at"}).
		AddRow(expectedCategories[0].ID, expectedCategories[0].UserID, expectedCategories[0].Name, expectedCategories[0].Description, expectedCategories[0].Icon, expectedCategories[0].CreatedAt, expectedCategories[0].UpdatedAt).
		AddRow(expectedCategories[1].ID, expectedCategories[1].UserID, expectedCategories[1].Name, expectedCategories[1].Description, expectedCategories[1].Icon, expectedCategories[1].CreatedAt, expectedCategories[1].UpdatedAt)

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnRows(rows)

	categories, err := GetPagedCategories(ctx, 2, itemsPerPage, 1, mockDBService)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, categories)
}
func TestInsertCategoryWithDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}

	expectedCategory := &Category{
		UserID:      1,
		Name:        "Test Category",
		Description: "Description",
		Icon:        "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectExec("^INSERT INTO categories (.+)").
		WithArgs(expectedCategory.UserID, expectedCategory.Name, expectedCategory.Description, expectedCategory.Icon, expectedCategory.CreatedAt, expectedCategory.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = InsertCategory(ctx, expectedCategory, mockDBService)
	assert.NoError(t, err)
}

func TestUpdateCategoryWithDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}

	existingCategory := &Category{
		ID:          1,
		UserID:      1,
		Name:        "Existing Category",
		Description: "Existing Description",
		Icon:        "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updatedCategory := &Category{
		ID:          1,
		UserID:      1,
		Name:        "Updated Category",
		Description: "Updated Description",
		Icon:        "",
		CreatedAt:   existingCategory.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	mock.ExpectExec("^UPDATE categories (.+) WHERE").
		WithArgs(updatedCategory.Name, updatedCategory.Description, updatedCategory.UpdatedAt, existingCategory.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(existingCategory.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "icon", "created_at", "updated_at"}).
			AddRow(updatedCategory.ID, updatedCategory.UserID, updatedCategory.Name, updatedCategory.Description, updatedCategory.Icon, updatedCategory.CreatedAt, updatedCategory.UpdatedAt))

	err = UpdateCategory(ctx, updatedCategory, mockDBService)
	assert.NoError(t, err)

	category, err := GetCategoryByID(ctx, updatedCategory.ID, updatedCategory.UserID, mockDBService)
	assert.NoError(t, err)
	assert.Equal(t, updatedCategory, category)
}

func TestGetAllCategoriesWithDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}

	userID := int64(1)

	expectedCategories := []*Category{
		{
			ID:          1,
			UserID:      userID,
			Name:        "Category 1",
			Description: "Description 1",
			Icon:        "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			UserID:      userID,
			Name:        "Category 2",
			Description: "Description 2",
			Icon:        "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "icon", "created_at", "updated_at"}).
		AddRow(expectedCategories[0].ID, expectedCategories[0].UserID, expectedCategories[0].Name, expectedCategories[0].Description, expectedCategories[0].Icon, expectedCategories[0].CreatedAt, expectedCategories[0].UpdatedAt).
		AddRow(expectedCategories[1].ID, expectedCategories[1].UserID, expectedCategories[1].Name, expectedCategories[1].Description, expectedCategories[1].Icon, expectedCategories[1].CreatedAt, expectedCategories[1].UpdatedAt)

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(userID).
		WillReturnRows(rows)

	categories, err := GetAllCategories(ctx, userID, mockDBService)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, categories)
}

func TestGetPagedCategoriesWithDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}

	userID := int64(1)
	page := 1
	itemsPerPage := 2

	expectedCategories := []*Category{
		{ID: 1, UserID: userID, Name: "Category 1", Description: "Description 1", Icon: "", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, UserID: userID, Name: "Category 2", Description: "Description 2", Icon: "", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "icon", "created_at", "updated_at"})
	for _, category := range expectedCategories {
		rows.AddRow(category.ID, category.UserID, category.Name, category.Description, category.Icon, category.CreatedAt, category.UpdatedAt)
	}

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(userID, (page-1)*itemsPerPage, itemsPerPage).
		WillReturnRows(rows)

	categories, err := GetPagedCategories(ctx, page, itemsPerPage, userID, mockDBService)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, categories)
	assert.Equal(t, 2, len(categories))
}



*/
