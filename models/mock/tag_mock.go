package mock

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"

	"github.com/stretchr/testify/mock"
)

// MockTagModel is a mock implementation of the TagService interface for testing
type MockTagModel struct {
	mock.Mock
}

// InsertTag mocks the InsertTag method
func (m *MockTagModel) InsertTag(ctx context.Context, tag *interfaces.Tag, otx ...*sql.Tx) error {
	args := m.Called(ctx, tag, otx)
	return args.Error(0)
}

// UpdateTag mocks the UpdateTag method
func (m *MockTagModel) UpdateTag(ctx context.Context, tag *interfaces.Tag, otx ...*sql.Tx) error {
	args := m.Called(ctx, tag, otx)
	return args.Error(0)
}

// Mock implementation of DeleteTagNew
func (m *MockTagModel) DeleteTag(ctx context.Context, tagID int64, scopes []int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, tagID, scopes, otx)
	return args.Error(0)
}

// Mock implementation of GetTagByIDNew
func (m *MockTagModel) GetTagByID(ctx context.Context, tagID int64, scopes []int64, otx ...*sql.Tx) (*interfaces.Tag, error) {
	args := m.Called(ctx, tagID, scopes, otx)
	return args.Get(0).(*interfaces.Tag), args.Error(1)
}

// Mock implementation of GetScopedTags
func (m *MockTagModel) GetScopedTags(ctx context.Context, scopes []int64, pagination interfaces.PaginationParams, otx ...*sql.Tx) ([]interfaces.Tag, error) {
	args := m.Called(ctx, scopes, pagination, otx)
	return args.Get(0).([]interfaces.Tag), args.Error(1)
}

// Mock implementation of GetTagByNameNew
func (m *MockTagModel) GetTagByName(ctx context.Context, name string, scopes []int64, otx ...*sql.Tx) (*interfaces.Tag, error) {
	args := m.Called(ctx, name, scopes, otx)
	return args.Get(0).(*interfaces.Tag), args.Error(1)
}

// Idiomatic interface compliance check.
// Ensure TagModel implements TagService
var _ interfaces.TagService = &MockTagModel{}
