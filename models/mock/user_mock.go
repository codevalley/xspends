package mock

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"

	"github.com/stretchr/testify/mock"
)

// MockUserModel is a mock implementation of the UserService interface for testing
type MockUserModel struct {
	mock.Mock
}

// InsertUser mocks the InsertUser method
func (m *MockUserModel) InsertUser(ctx context.Context, user *interfaces.User, otx ...*sql.Tx) error {
	args := m.Called(ctx, user, otx)
	return args.Error(0)
}

// UpdateUser mocks the UpdateUser method
func (m *MockUserModel) UpdateUser(ctx context.Context, user *interfaces.User, otx ...*sql.Tx) error {
	args := m.Called(ctx, user, otx)
	return args.Error(0)
}

// DeleteUser mocks the DeleteUser method
func (m *MockUserModel) DeleteUser(ctx context.Context, id int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, id, otx)
	return args.Error(0)
}

// GetUserByID mocks the GetUserByID method
func (m *MockUserModel) GetUserByID(ctx context.Context, id int64, otx ...*sql.Tx) (*interfaces.User, error) {
	args := m.Called(ctx, id, otx)
	return args.Get(0).(*interfaces.User), args.Error(1)
}

// GetUserByUsername mocks the GetUserByUsername method
func (m *MockUserModel) GetUserByUsername(ctx context.Context, username string, otx ...*sql.Tx) (*interfaces.User, error) {
	args := m.Called(ctx, username, otx)
	return args.Get(0).(*interfaces.User), args.Error(1)
}

// UserExists mocks the UserExists method
func (m *MockUserModel) UserExists(ctx context.Context, username, email string, otx ...*sql.Tx) (bool, error) {
	args := m.Called(ctx, username, email, otx)
	return args.Bool(0), args.Error(1)
}

// UserIDExists mocks the UserIDExists method
func (m *MockUserModel) UserIDExists(ctx context.Context, id int64, otx ...*sql.Tx) (bool, error) {
	args := m.Called(ctx, id, otx)
	return args.Bool(0), args.Error(1)
}

// Idiomatic interface compliance check.
// Ensure UserModel implements UserService
var _ interfaces.UserService = &MockUserModel{}
