package mock

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockGroupModel struct {
	mock.Mock
}

// Implementing CreateGroup method of GroupService interface
func (m *MockGroupModel) CreateGroup(ctx context.Context, group *interfaces.Group, userIDs []int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, group, userIDs, otx)
	return args.Error(0)
}
func (m *MockGroupModel) UpdateGroup(ctx context.Context, group *interfaces.Group, requestingUserID int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, group, requestingUserID, otx)
	return args.Error(0)
}

// Implementing DeleteGroup method of GroupService interface
func (m *MockGroupModel) DeleteGroup(ctx context.Context, groupID int64, requestingUserID int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, groupID, requestingUserID, otx)
	return args.Error(0)
}

// Implementing GetGroupByID method of GroupService interface
func (m *MockGroupModel) GetGroupByID(ctx context.Context, groupID int64, requestingUserID int64, otx ...*sql.Tx) (*interfaces.Group, error) {
	args := m.Called(ctx, groupID, requestingUserID, otx)
	return args.Get(0).(*interfaces.Group), args.Error(1)
}

// Implementing GetGroupByScope method of GroupService interface
func (m *MockGroupModel) GetGroupByScope(ctx context.Context, scopeID int64, requestingUserID int64, otx ...*sql.Tx) (*interfaces.Group, error) {
	args := m.Called(ctx, scopeID, requestingUserID, otx)
	return args.Get(0).(*interfaces.Group), args.Error(1)
}

// Ensure MockGroupModel implements GroupService interface
var _ interfaces.GroupService = &MockGroupModel{}
