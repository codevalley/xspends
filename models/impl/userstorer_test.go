package impl

import (
	"context"
	"testing"
	"time"
	"xspends/models/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserStorer_Load(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()
	// Set up SQL Mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Assuming a specific user for testing purposes
	username := "testuser"
	expectedUser := &interfaces.User{
		Username: username,
		Email:    "test@example.com",
	}

	// Prepare the mock response
	rows := sqlmock.NewRows([]string{"id", "username", "name", "email", "currency", "password"}).
		AddRow("1", expectedUser.Username, "Test User", expectedUser.Email, "USD", "hashedpassword")

	// Set up the expected SQL query that will be run
	// Note that the sqlquery variable comes from your actual GetUserByUsername method,
	// so ensure it matches exactly with what's being executed there.
	sqlquery := "SELECT id, username, name, email, currency, password FROM users WHERE username = ?"
	mock.ExpectQuery(sqlquery).
		WithArgs(username).
		WillReturnRows(rows)

	// Inject the mock executor into your UserModel or database service as needed.
	// Assuming ModelsService is where the UserModel exists and it has a method to set the executor.
	mockModelService.DBService.Executor = db
	userStorer := NewUserStorer()

	// Call the Load method which internally calls GetUserByUsername
	user, err := userStorer.Load(context.Background(), username)
	assert.NoError(t, err)

	// Convert the returned authboss.User to your *interfaces.User
	loadedUser, ok := user.(*interfaces.User)
	assert.True(t, ok)
	assert.Equal(t, expectedUser.Username, loadedUser.Username)
	assert.Equal(t, expectedUser.Email, loadedUser.Email)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserStorer_Save(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	// Mock setup
	testUser := &interfaces.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Currency:  "USD",
		UpdatedAt: time.Now(), // This will be dynamic, consider using a matcher if needed
	}

	// This should match the actual SQL query string
	expectedSQL := "UPDATE users SET currency = ?, email = ?, name = ?, password = ?, updated_at = ?, username = ? WHERE id = ?"

	// Mock the database call
	mockExecutor.EXPECT().
		ExecContext(
			gomock.Any(), // The context
			expectedSQL,  // The SQL query
			gomock.Any(), // Use gomock.Any() for dynamic values or specify each expected value
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(sqlmock.NewResult(1, 1), nil) // Assuming successful execution

	userStorer := NewUserStorer()
	err = userStorer.Save(context.Background(), testUser)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserStorer_LoadByConfirmSelector(t *testing.T) {
	userStorer := NewUserStorer()
	_, err := userStorer.LoadByConfirmSelector(context.Background(), "selector")
	assert.Error(t, err)
}
func TestUserStorer_LoadByRecoverSelector(t *testing.T) {
	userStorer := NewUserStorer()
	_, err := userStorer.LoadByRecoverSelector(context.Background(), "selector")
	assert.Error(t, err)
}
