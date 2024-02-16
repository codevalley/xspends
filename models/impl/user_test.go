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
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.UserModel = NewUserModel()
		config.ScopeModel = NewScopeModel()
		config.UserScopeModel = NewUserScopeModel()
	})
	defer tearDown()

	user := &interfaces.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Ensure the SQL query and other arguments match exactly with those used in the InsertUser method
	mockExecutor.EXPECT().
		ExecContext(
			gomock.Any(), // The context
			"INSERT INTO scopes (scope_id,type) VALUES (?,?)", // The SQL query
			gomock.Any(), // Match each argument
			gomock.Any(),
		).                                   //WithArgs(sqlmock.AnyArg(), "user").
		Return(sqlmock.NewResult(1, 1), nil) // Simulate successful execution

	// Ensure the SQL query and other arguments match exactly with those used in the InsertUser method
	mockExecutor.EXPECT().
		ExecContext(
			gomock.Any(), // The context
			"INSERT INTO user_scopes (user_id,scope_id,role) VALUES (?,?,?) ON DUPLICATE KEY UPDATE role = VALUES(role)", // The SQL query
			gomock.Any(), // Match each argument
			gomock.Any(),
			gomock.Any(),
		).                                   //.WithArgs(newUser.ID, scopeID, "owner")
		Return(sqlmock.NewResult(1, 1), nil) // Simulate successful execution
	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)
	err := ModelsService.UserModel.InsertUser(ctx, user)
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.UserModel = NewUserModel()
	})
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

	err := ModelsService.UserModel.UpdateUser(ctx, user)
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.UserModel = NewUserModel()
		config.ScopeModel = NewScopeModel()
		config.UserScopeModel = NewUserScopeModel()
	})
	defer tearDown()

	userID := int64(1)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDBService := &DBService{Executor: db}
	mockModelService := &ModelsServiceContainer{
		DBService:      mockDBService,
		UserModel:      NewUserModel(),
		ScopeModel:     NewScopeModel(),
		UserScopeModel: NewUserScopeModel(),
		// Initialize other models as necessary
	}
	ModelsService = mockModelService

	rows := sqlmock.NewRows([]string{"user_id", "username", "name", "email", "scope_id", "currency", "password"}).
		AddRow(userID, "testuser", "Test User", "test@example.com", 1, "USD", "hashedpassword")
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(userID).WillReturnRows(rows)

	mock.ExpectExec("^DELETE FROM users WHERE").WithArgs(userID).WillReturnResult(sqlmock.NewResult(1, 1))
	err = ModelsService.UserModel.DeleteUser(ctx, userID)
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
	mockModelService := &ModelsServiceContainer{
		DBService: mockDBService,
		UserModel: NewUserModel(),
		// Initialize other models as necessary
	}
	ModelsService = mockModelService

	// Set up expectations
	userID := int64(1)
	rows := sqlmock.NewRows([]string{"user_id", "username", "name", "email", "scope_id", "currency", "password"}).
		AddRow(userID, "testuser", "Test User", "test@example.com", 1, "USD", "hashedpassword")
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(userID).WillReturnRows(rows)

	// Call the function under test
	ctx := context.Background()
	user, err := ModelsService.UserModel.GetUserByID(ctx, userID)

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
	mockModelService := &ModelsServiceContainer{
		DBService: mockDBService,
		UserModel: NewUserModel(),
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

	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "scope", "currency", "password"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Name, expectedUser.Email, expectedUser.Scope, expectedUser.Currency, expectedUser.Password)
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(username).WillReturnRows(rows)

	// Call the function under test
	user, err := ModelsService.UserModel.GetUserByUsername(ctx, username)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Scope, user.Scope)
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
	mockModelService := &ModelsServiceContainer{
		DBService: mockDBService,
		UserModel: NewUserModel(),
		// Initialize other models as necessary
	}
	ModelsService = mockModelService

	// Set up expectations
	username := "existinguser"
	email := "existing@example.com"

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(username, email).WillReturnRows(rows)

	// Call the function under test
	exists, err := ModelsService.UserModel.UserExists(ctx, username, email)

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
	mockModelService := &ModelsServiceContainer{
		DBService: mockDBService,
		UserModel: NewUserModel(),
		// Initialize other models as necessary
	}
	ModelsService = mockModelService

	// Set up expectations
	userID := int64(1)

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(userID).WillReturnRows(rows)

	// Call the function under test
	exists, err := ModelsService.UserModel.UserIDExists(ctx, userID)

	// Assertions
	assert.NoError(t, err)
	assert.True(t, exists)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserIDExists_UserIDNotFound(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDBService := &DBService{Executor: db}
	mockModelService := &ModelsServiceContainer{
		DBService: mockDBService,
		UserModel: NewUserModel(),
		// Initialize other models as necessary
	}
	ModelsService = mockModelService
	userID := int64(1) // Non-existent userID
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(userID).WillReturnRows(sqlmock.NewRows([]string{"exists"}))

	// Call the function under test
	exists, err := ModelsService.UserModel.UserIDExists(ctx, userID)

	// Assertions
	// 1. Ensure no error is returned (indicating a successful query)
	assert.NoError(t, err)

	// 2. Assert that 'exists' is false, indicating the userID was not found
	assert.False(t, exists)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPutPID(t *testing.T) {
	var user interfaces.User
	pid := "123"
	user.PutPID(pid)
	if user.ID != 123 {
		t.Errorf("PutPID was incorrect, got: %d, want: %d.", user.ID, 123)
	}
}

func TestGetPID(t *testing.T) {
	user := interfaces.User{ID: 123}
	got := user.GetPID()
	want := "123"
	if got != want {
		t.Errorf("GetPID was incorrect, got: %s, want: %s.", got, want)
	}
}

func TestPutPassword(t *testing.T) {
	var user interfaces.User
	password := "securepassword"
	user.PutPassword(password)
	if user.Password != "securepassword" {
		t.Errorf("PutPassword was incorrect, got: %s, want: %s.", user.Password, "securepassword")
	}
}

// Add more test cases for UpdateUser, DeleteUser, GetUserByUsername, UserExists, UserIDExists
