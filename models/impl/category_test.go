package impl

import (
	"testing"
	"time"
	"xspends/models/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// Insert your test functions here, following the pattern demonstrated below.

func TestInsertCategoryWithInvalidInput(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	invalidCategories := []*interfaces.Category{
		{UserID: 0, Name: "Test Category", Description: "Description"},
		{UserID: 1, Name: "", Description: "Description"},
		{UserID: 1, Name: "Test Category", Description: string(make([]byte, 501))},
	}

	for _, invalidCategory := range invalidCategories {
		err := ModelsService.CategoryModel.InsertCategory(ctx, invalidCategory)
		assert.EqualError(t, err, ErrInvalidInput)
	}
}

// TestInsertCategoryWithDatabaseError verifies that the function returns an error for database errors.
func TestInsertCategoryWithDatabaseError(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	category := &interfaces.Category{UserID: 1, Name: "Test Category", Description: "Description"}

	err := ModelsService.CategoryModel.InsertCategory(ctx, category)
	assert.EqualError(t, err, "executing insert statement failed: database error")
}

// TestUpdateCategoryWithDatabaseError verifies that the function returns an error for database errors.
func TestUpdateCategoryWithDatabaseError(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	category := &interfaces.Category{ID: 1, UserID: 1, Name: "Updated Category", Description: "Updated Description"}
	err := ModelsService.CategoryModel.UpdateCategory(ctx, category)
	assert.EqualError(t, err, "executing update statement failed: database error")
}

// TestDeleteCategoryWithDatabaseError verifies that the function returns an error for database errors.
func TestDeleteCategoryWithDatabaseError(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	mockExecutor.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

	err := ModelsService.CategoryModel.DeleteCategory(ctx, 1, 1)
	assert.EqualError(t, err, "executing delete statement failed: database error")
}

// TestGetAllCategoriesWithDatabaseError verifies that the function returns an error for database errors.
func TestGetAllCategoriesWithDatabaseError(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnError(errors.New("database error"))

	_, err = ModelsService.CategoryModel.GetAllCategories(ctx, 1)
	assert.EqualError(t, err, "querying categories failed: database error")
}

// TestGetCategoryByIDWithCategoryNotFound verifies that the function returns an error for a non-existent category.
func TestGetCategoryByIDWithCategoryNotFound(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnRows(sqlmock.NewRows([]string{}))

	_, err = ModelsService.CategoryModel.GetCategoryByID(ctx, 1, 1)
	assert.EqualError(t, err, "category not found")
}

// TestGetCategoryByIDWithDatabaseError verifies that the function returns an error for database errors.
func TestGetCategoryByIDWithDatabaseError(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnError(errors.New("database error"))

	_, err = ModelsService.CategoryModel.GetCategoryByID(ctx, 1, 1)
	assert.EqualError(t, err, "querying category by ID failed: database error")
}

// TestGetPagedCategoriesWithDatabaseError verifies that the function returns an error for database errors.
func TestGetPagedCategoriesWithDatabaseError(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnError(errors.New("database error"))

	_, err = ModelsService.CategoryModel.GetPagedCategories(ctx, 1, 10, 1)
	assert.EqualError(t, err, "querying paginated categories failed: database error")
}

// TestCategoryIDExistsWithDatabaseError verifies that the function returns an error for database errors.
func TestCategoryIDExistsWithDatabaseError(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnError(errors.New("database error"))

	exists, err := ModelsService.CategoryModel.CategoryIDExists(ctx, 1, 1, nil)
	assert.EqualError(t, err, "checking category existence failed: database error")
	assert.False(t, exists)
}

// TestGetCategoryByIDWithEmptyIcon verifies that the category is returned with an empty icon when not set.
func TestGetCategoryByIDWithEmptyIcon(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	expectedCategory := &interfaces.Category{
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

	category, err := ModelsService.CategoryModel.GetCategoryByID(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, category)
}

// TestGetAllCategoriesWithEmptyResults verifies that an empty slice is returned when there are no categories.
func TestGetAllCategoriesWithEmptyResults(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	rows := sqlmock.NewRows([]string{})

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnRows(rows)

	categories, err := ModelsService.CategoryModel.GetAllCategories(ctx, 1)
	assert.NoError(t, err)
	assert.Empty(t, categories)
}

// TestCategoryIDExistsWithNonExistentCategory verifies that false is returned for non-existent category.
func TestCategoryIDExistsWithNonExistentCategory(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	categoryID := int64(1)
	userID := int64(1)

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	actualExists, err := ModelsService.CategoryModel.CategoryIDExists(ctx, categoryID, userID)
	assert.NoError(t, err)
	assert.False(t, actualExists)
}

// TestDeleteCategoryWithDatabase verifies that DeleteCategory works with a mock database.
func TestDeleteCategoryWithDatabase(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	categoryID := int64(1)
	userID := int64(1)

	mock.ExpectExec("^DELETE FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = ModelsService.CategoryModel.DeleteCategory(ctx, categoryID, userID)
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
	mockModelService := &ModelsServiceContainer{
		DBService:     mockDBService,
		CategoryModel: &CategoryModel{},
	}
	ModelsService = mockModelService

	categoryID := int64(1)
	userID := int64(1)

	expectedCategory := &interfaces.Category{
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

	category, err := ModelsService.CategoryModel.GetCategoryByID(ctx, categoryID, userID)
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
	mockModelService := &ModelsServiceContainer{
		DBService:     mockDBService,
		CategoryModel: &CategoryModel{},
	}
	ModelsService = mockModelService

	categoryID := int64(1)
	userID := int64(1)

	exists := true

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").
		WithArgs(categoryID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(categoryID))

	actualExists, err := ModelsService.CategoryModel.CategoryIDExists(ctx, categoryID, userID)
	assert.NoError(t, err)
	assert.Equal(t, exists, actualExists)
}

// TestGetPagedCategoriesWithEmptyResults verifies that an empty slice is returned when there are no categories for pagination.
func TestGetPagedCategoriesWithEmptyResults(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	rows := sqlmock.NewRows([]string{})

	mock.ExpectQuery("^SELECT (.+) FROM categories WHERE").WillReturnRows(rows)

	categories, err := ModelsService.CategoryModel.GetPagedCategories(ctx, 1, 10, 1, nil)
	assert.NoError(t, err)
	assert.Empty(t, categories)
}

// TestUpdateCategoryWithValidInput verifies that UpdateCategory successfully updates a category.
func TestUpdateCategoryWithValidInput(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.CategoryModel = &CategoryModel{}
	})
	defer tearDown()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ModelsService.DBService.Executor = db

	category := &interfaces.Category{ID: 1, UserID: 1, Name: "Updated Category", Description: "Updated Description"}

	mock.ExpectExec("^UPDATE categories SET").WithArgs(category.Name, category.Description, category.Icon, sqlmock.AnyArg(), category.ID, category.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = ModelsService.CategoryModel.UpdateCategory(ctx, category)
	assert.NoError(t, err)
}

// Add any additional test cases as required for your application...
