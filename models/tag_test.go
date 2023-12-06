package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Insert a valid tag with a name and user ID
func TestInsertValidTag(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	tag := &Tag{
		UserID:    1,
		Name:      "Test Tag",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := InsertTag(ctx, tag, mockDBService)
	assert.NoError(t, err)
}

// Update an existing tag with a valid name and user ID
func TestUpdateExistingTag(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	tag := &Tag{
		ID:        1,
		UserID:    1,
		Name:      "Updated Tag",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := UpdateTag(ctx, tag, mockDBService)
	assert.NoError(t, err)
}

// Delete an existing tag with a valid tag ID and user ID
func TestDeleteExistingTag(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	tagID := int64(1)
	userID := int64(1)

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := DeleteTag(ctx, tagID, userID, mockDBService)
	assert.NoError(t, err)
}

// Retrieve an existing tag by ID with a valid tag ID and user ID
func TestRetrieveExistingTagByID(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	tagID := int64(1)
	userID := int64(1)

	mockRows := sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
		AddRow(tagID, userID, "Test Tag", time.Now(), time.Now())

	mockExecutor.EXPECT().
		QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockRows).
		Times(1)

	tag, err := GetTagByID(ctx, tagID, userID, mockDBService)
	assert.NoError(t, err)
	assert.NotNil(t, tag)
}

// Retrieve all tags for a user with a valid user ID and pagination parameters
func TestRetrieveAllTagsForUser(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	userID := int64(1)
	pagination := PaginationParams{
		Limit:  10,
		Offset: 0,
	}

	mockRows := sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
		AddRow(1, userID, "Tag 1", time.Now(), time.Now()).
		AddRow(2, userID, "Tag 2", time.Now(), time.Now())

	mockExecutor.EXPECT().
		QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockRows, nil).
		Times(1)

	tags, err := GetAllTags(ctx, userID, pagination, mockDBService)
	assert.NoError(t, err)
	assert.NotNil(t, tags)
	assert.Equal(t, 2, len(tags))
}

// Insert a tag with an invalid user ID or name
func TestInsertInvalidTag(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	tag := &Tag{
		UserID:    0,
		Name:      "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := InsertTag(ctx, tag, mockDBService)
	assert.Error(t, err)
}

func TestUpdateInvalidTag(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	tag := &Tag{
		ID:        1,
		UserID:    0,
		Name:      "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := UpdateTag(ctx, tag, mockDBService)
	assert.NotNil(t, err, "Expected error due to invalid tag")
	assert.Contains(t, err.Error(), "invalid input for tag", "Expected 'invalid input for tag' error")
}
