package impl

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"xspends/models/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Insert a valid tag with a name and user ID
func TestInsertValidTag(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TagModel = NewTagModel()
	})
	defer tearDown()

	tag := &interfaces.Tag{
		UserID:    1,
		ScopeID:   1,
		Name:      "Test Tag",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := ModelsService.TagModel.InsertTag(ctx, tag)
	assert.NoError(t, err)
}

// Update an existing tag with a valid name and user ID
func TestUpdateExistingTag(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TagModel = NewTagModel()
	})
	defer tearDown()

	tag := &interfaces.Tag{
		ID:        1,
		UserID:    1,
		ScopeID:   1,
		Name:      "Updated Tag",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := ModelsService.TagModel.UpdateTag(ctx, tag)
	assert.NoError(t, err)
}

// Delete an existing tag with a valid tag ID and user ID
func TestDeleteExistingTag(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TagModel = NewTagModel()
	})
	defer tearDown()

	tagID := int64(1)
	scopes := []int64{1}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := ModelsService.TagModel.DeleteTag(ctx, tagID, scopes)
	assert.NoError(t, err)
}

// Insert a tag with an invalid user ID or name
func TestInsertInvalidTag(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TagModel = NewTagModel()
	})
	defer tearDown()

	tag := &interfaces.Tag{
		UserID:    0,
		Name:      "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ModelsService.TagModel.InsertTag(ctx, tag)
	assert.Error(t, err)
}

func TestUpdateInvalidTag(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TagModel = NewTagModel()
	})
	defer tearDown()

	tag := &interfaces.Tag{
		ID:        1,
		UserID:    0,
		Name:      "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := ModelsService.TagModel.UpdateTag(ctx, tag)
	assert.NotNil(t, err, "Expected error due to invalid tag")
	assert.Contains(t, err.Error(), "invalid input for tag", "Expected 'invalid input for tag' error")
}

// Retrieve an existing tag by ID with a valid tag ID and user ID
func TestRetrieveExistingTagByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDBService := &DBService{Executor: db}
	mockModelService := &ModelsServiceContainer{
		DBService: mockDBService,
		TagModel:  NewTagModel(),
	}
	ModelsService = mockModelService

	tagID := int64(1)
	userID := int64(1)
	scopes := []int64{1, 2, 3} // Example scopes

	// Adjusted mockRows to match expected columns order and data
	mockRows := sqlmock.NewRows([]string{"tag_id", "user_id", "name", "scope_id", "created_at", "updated_at"}).
		AddRow(tagID, userID, "Test Tag", scopes[0], time.Now(), time.Now())

		// Adjusting regex to match the actual query pattern and placeholders
	mock.ExpectQuery(`SELECT tag_id, user_id, name, scope_id, created_at, updated_at FROM tags WHERE scope_id IN \(\?,\?,\?\) AND tag_id = \?`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tagID). // Using AnyArg for scopes and explicitly for tagID
		WillReturnRows(mockRows)

	ctx := context.Background()
	tag, err := ModelsService.TagModel.GetTagByID(ctx, tagID, scopes)

	assert.NoError(t, err)
	assert.NotNil(t, tag)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Retrieve all tags for a user with a valid user ID and pagination parameters
// Retrieve all tags for a user with a valid user ID and pagination parameters
func TestRetrieveAllTagsForUser(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TagModel = NewTagModel()
	})
	defer tearDown()

	userID := int64(1)
	scopes := []int64{1, 2, 3} // Example scopes
	pagination := interfaces.PaginationParams{
		Limit:  10,
		Offset: 0,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}
	mockModelService := &ModelsServiceContainer{
		DBService: mockDBService,
		TagModel:  NewTagModel(),
	}
	ModelsService = mockModelService

	mockRows := sqlmock.NewRows([]string{"tag_id", "user_id", "name", "scope_id", "created_at", "updated_at"}).
		AddRow(1, userID, "Tag 1", scopes[0], time.Now(), time.Now()).
		AddRow(2, userID, "Tag 2", scopes[0], time.Now(), time.Now())

		// Use a lenient regular expression that focuses on key parts of the query
		// Update ExpectQuery to reflect the actual arguments used in GetAllTags
	mock.ExpectQuery(`SELECT tag_id, user_id, name, scope_id, created_at, updated_at FROM tags WHERE scope_id IN \(\?,\?,\?\) LIMIT 10 OFFSET 0`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()). // Only include userID if that's the only argument used
		WillReturnRows(mockRows)

	ctx := context.Background()
	tags, err := ModelsService.TagModel.GetScopedTags(ctx, scopes, pagination)

	assert.NoError(t, err)
	assert.NotNil(t, tags)
	assert.Equal(t, 2, len(tags))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTagByName(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.TagModel = NewTagModel()
	})
	defer tearDown()

	userID := int64(1)
	scopes := []int64{1, 2, 3} // Example scopes
	tagName := "Test Tag"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}
	mockModelService := &ModelsServiceContainer{
		DBService: mockDBService,
		TagModel:  NewTagModel(),
	}
	ModelsService = mockModelService

	mockRows := sqlmock.NewRows([]string{"tag_id", "user_id", "name", "scope_id", "created_at", "updated_at"}).
		AddRow(1, userID, tagName, scopes[0], time.Now(), time.Now())

	// Set up the expected query with sqlmock
	mock.ExpectQuery(`SELECT tag_id, user_id, name, scope_id, created_at, updated_at FROM tags WHERE name = \? AND scope_id IN \(\?,\?,\?\)`).
		WithArgs(tagName, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(mockRows)

	ctx := context.Background()
	tag, err := ModelsService.TagModel.GetTagByName(ctx, tagName, scopes)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, tagName, tag.Name)
	assert.Equal(t, userID, tag.UserID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
