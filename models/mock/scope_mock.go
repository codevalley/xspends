package mock

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockScopeModel struct {
	mock.Mock
}

var _ interfaces.ScopeService = &MockScopeModel{}

func (m *MockScopeModel) CreateScope(ctx context.Context, scopeType string, otx ...*sql.Tx) (int64, error) {
	args := m.Called(ctx, scopeType, otx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockScopeModel) GetScope(ctx context.Context, scopeID int64, otx ...*sql.Tx) (*interfaces.Scope, error) {
	args := m.Called(ctx, scopeID, otx)
	return args.Get(0).(*interfaces.Scope), args.Error(1)
}

func (m *MockScopeModel) DeleteScope(ctx context.Context, scopeID int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, scopeID, otx)
	return args.Error(0)
}

func (m *MockScopeModel) ScopeIDExists(ctx context.Context, scopeId int64, otx ...*sql.Tx) (bool, error) {
	args := m.Called(ctx, scopeId, otx)
	return args.Bool(0), args.Error(1)
}
