package impl

import (
	"database/sql"
	"testing"
	"time"
	"xspends/models/interfaces"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	user := &interfaces.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)
	err := mockModelService.UserModel.InsertUser(ctx, user)
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	user := &interfaces.User{
		ID:        1,
		Username:  "updateduser",
		Email:     "updated@example.com",
		Password:  "newpassword123",
		UpdatedAt: time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := mockModelService.UserModel.UpdateUser(ctx, user)
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	userID := int64(1)

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := mockModelService.UserModel.DeleteUser(ctx, userID)
	assert.NoError(t, err)
}

// Add more test cases for GetUserByID, GetUserByUsername, UserExists, UserIDExists

// Add more test cases for UpdateUser, DeleteUser, GetUserByUsername, UserExists, UserIDExists
