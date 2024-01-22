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

// DeleteGroup implements interfaces.GroupService.
func (*MockGroupModel) DeleteGroup(ctx context.Context, groupID int64, otx ...*sql.Tx) error {
	panic("unimplemented")
}

// GetGroupByID implements interfaces.GroupService.
func (*MockGroupModel) GetGroupByID(ctx context.Context, groupID int64, otx ...*sql.Tx) (*interfaces.Group, error) {
	panic("unimplemented")
}

// InsertGroup implements interfaces.GroupService.
func (*MockGroupModel) InsertGroup(ctx context.Context, group *interfaces.Group, otx ...*sql.Tx) error {
	panic("unimplemented")
}

// UpdateGroup implements interfaces.GroupService.
func (*MockGroupModel) UpdateGroup(ctx context.Context, group *interfaces.Group, otx ...*sql.Tx) error {
	panic("unimplemented")
}

var _ interfaces.GroupService = &MockGroupModel{}

// Implement methods like InsertGroup, UpdateGroup, etc., similar to the MockCategoryModel methods.
