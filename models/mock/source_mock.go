package mock

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"

	"github.com/stretchr/testify/mock"
)

// MockSourceModel is a mock implementation of SourceService.
type MockSourceModel struct {
	mock.Mock
}

// Ensure MockSourceModel implements SourceService.
var _ interfaces.SourceService = &MockSourceModel{}

func (m *MockSourceModel) InsertSource(ctx context.Context, source *interfaces.Source, otx ...*sql.Tx) error {
	args := m.Called(ctx, source, otx)
	return args.Error(0)
}

func (m *MockSourceModel) UpdateSource(ctx context.Context, source *interfaces.Source, otx ...*sql.Tx) error {
	args := m.Called(ctx, source, otx)
	return args.Error(0)
}

func (m *MockSourceModel) DeleteSource(ctx context.Context, sourceID int64, userID int64, otx ...*sql.Tx) error {
	args := m.Called(ctx, sourceID, userID, otx)
	return args.Error(0)
}

func (m *MockSourceModel) GetSourceByID(ctx context.Context, sourceID int64, userID int64, otx ...*sql.Tx) (*interfaces.Source, error) {
	args := m.Called(ctx, sourceID, userID, otx)
	return args.Get(0).(*interfaces.Source), args.Error(1)
}

func (m *MockSourceModel) GetSources(ctx context.Context, userID int64, otx ...*sql.Tx) ([]interfaces.Source, error) {
	args := m.Called(ctx, userID, otx)
	return args.Get(0).([]interfaces.Source), args.Error(1)
}

func (m *MockSourceModel) SourceIDExists(ctx context.Context, sourceID int64, userID int64, otx ...*sql.Tx) (bool, error) {
	args := m.Called(ctx, sourceID, userID, otx)
	return args.Bool(0), args.Error(1)
}

// Idiomatic interface compliance check.
// Ensure SourceModel implements SourceService
var _ interfaces.SourceService = &MockSourceModel{}
