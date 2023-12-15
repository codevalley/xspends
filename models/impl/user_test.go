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
func TestGetUserByID(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDBService := &DBService{Executor: db}
	mockModelService = &ModelsServiceContainer{
		DBService: mockDBService,
		UserModel: &UserModel{},
		// Initialize other models as necessary
	}
	ModelsService = mockModelService

	// Set up expectations
	userID := int64(1)
	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "currency", "password"}).
		AddRow(userID, "testuser", "Test User", "test@example.com", "USD", "hashedpassword")
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(userID).WillReturnRows(rows)

	// Call the function under test
	ctx := context.Background()
	user, err := mockModelService.UserModel.GetUserByID(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "USD", user.Currency)
	assert.Equal(t, "hashedpassword", user.Password)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByUsername(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDBService := &DBService{Executor: db}
	mockModelService = &ModelsServiceContainer{
		DBService: mockDBService,
		UserModel: &UserModel{},
		// Initialize other models as necessary
	}
	ModelsService = mockModelService

	// Set up expectations
	username := "testuser"
	expectedUser := &interfaces.User{
		ID:       1,
		Username: username,
		Name:     "Test User",
		Email:    "test@example.com",
		Currency: "USD",
		Password: "hashedpassword",
	}

	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "currency", "password"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Name, expectedUser.Email, expectedUser.Currency, expectedUser.Password)
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(username).WillReturnRows(rows)

	// Call the function under test
	user, err := mockModelService.UserModel.GetUserByUsername(ctx, username)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Currency, user.Currency)
	assert.Equal(t, expectedUser.Password, user.Password)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestUserExists(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDBService := &DBService{Executor: db}
	mockModelService = &ModelsServiceContainer{
		DBService: mockDBService,
		UserModel: &UserModel{},
		// Initialize other models as necessary
	}
	ModelsService = mockModelService

	// Set up expectations
	username := "existinguser"
	email := "existing@example.com"

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(username, email).WillReturnRows(rows)

	// Call the function under test
	exists, err := mockModelService.UserModel.UserExists(ctx, username, email)

	// Assertions
	assert.NoError(t, err)
	assert.True(t, exists)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestUserIDExists(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDBService := &DBService{Executor: db}
	mockModelService = &ModelsServiceContainer{
		DBService: mockDBService,
		UserModel: &UserModel{},
		// Initialize other models as necessary
	}
	ModelsService = mockModelService

	// Set up expectations
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(userID).WillReturnRows(rows)

	// Call the function under test
	exists, err := mockModelService.UserModel.UserIDExists(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	assert.True(t, exists)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Add more test cases for UpdateUser, DeleteUser, GetUserByUsername, UserExists, UserIDExists
