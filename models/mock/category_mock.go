package mock

import (
	context "context"
	sql "database/sql"
	"xspends/models/interfaces"

	"github.com/stretchr/testify/mock"
)

// MockCategoryModel is a mock implementation of CategoryService
type MockCategoryModel struct {
	mock.Mock
}

// Ensure MockCategoryModel implements CategoryService
var _ interfaces.CategoryService = &MockCategoryModel{}

func (m *MockCategoryModel) InsertCategory(ctx context.Context, category *interfaces.Category, otx ...*sql.Tx) error {
	args := m.Called(ctx, category, otx)
	return args.Error(0)
}

func (m *MockCategoryModel) UpdateCategory(ctx context.Context, category *interfaces.Category, otx ...*sql.Tx) error {
	args := m.Called(ctx, category, otx)
	return args.Error(0)
}

func (m *MockCategoryModel) GetAllCategories(ctx context.Context, userID int64, otx ...*sql.Tx) ([]interfaces.Category, error) {
	args := m.Called(ctx, userID, otx)
	return args.Get(0).([]interfaces.Category), args.Error(1)
}

func (m *MockCategoryModel) GetPagedCategories(ctx context.Context, page int, itemsPerPage int, userID int64, otx ...*sql.Tx) ([]interfaces.Category, error) {
	args := m.Called(ctx, page, itemsPerPage, userID, otx)
	return args.Get(0).([]interfaces.Category), args.Error(1)
}

func (m *MockCategoryModel) DeleteCategory(ctx context.Context, categoryID int64, scopes []int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, categoryID, scopes, otx)
	return args.Error(0)
}

func (m *MockCategoryModel) GetAllScopedCategories(ctx context.Context, scopes []int64, otx ...*sql.Tx) ([]interfaces.Category, error) {
	args := m.Called(ctx, scopes, otx)
	return args.Get(0).([]interfaces.Category), args.Error(1)
}

func (m *MockCategoryModel) GetCategoryByID(ctx context.Context, categoryID int64, scopes []int64, otx ...*sql.Tx) (*interfaces.Category, error) {
	args := m.Called(ctx, categoryID, scopes, otx)
	return args.Get(0).(*interfaces.Category), args.Error(1)
}

func (m *MockCategoryModel) GetScopedCategories(ctx context.Context, page int, itemsPerPage int, scopes []int64, otx ...*sql.Tx) ([]interfaces.Category, error) {
	args := m.Called(ctx, page, itemsPerPage, scopes, otx)
	return args.Get(0).([]interfaces.Category), args.Error(1)
}

func (m *MockCategoryModel) CategoryIDExists(ctx context.Context, categoryID int64, scopes []int64, otx ...*sql.Tx) (bool, error) {
	args := m.Called(ctx, categoryID, scopes, otx)
	return args.Bool(0), args.Error(1)
}

// Idiomatic interface compliance check.
// Ensure CategoryModel implements CategoryService
var _ interfaces.CategoryService = &MockCategoryModel{}
