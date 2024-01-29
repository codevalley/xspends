package mock

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockUserScopeModel struct {
	mock.Mock
}

var _ interfaces.UserScopeService = &MockUserScopeModel{}

// Mock implementation of UpsertUserScope
func (m *MockUserScopeModel) UpsertUserScope(ctx context.Context, userID, scopeID int64, role string, otx ...*sql.Tx) error {
	args := m.Called(ctx, userID, scopeID, role, otx)
	return args.Error(0)
}

// Mock implementation of GetUserScope
func (m *MockUserScopeModel) GetUserScope(ctx context.Context, userID, scopeID int64, otx ...*sql.Tx) (*interfaces.UserScope, error) {
	args := m.Called(ctx, userID, scopeID, otx)
	return args.Get(0).(*interfaces.UserScope), args.Error(1)
}

// Mock implementation of DeleteUserScope
func (m *MockUserScopeModel) DeleteUserScope(ctx context.Context, userID, scopeID int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, userID, scopeID, otx)
	return args.Error(0)
}
