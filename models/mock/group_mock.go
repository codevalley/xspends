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

func (*MockGroupModel) CreateGroup(ctx context.Context, group *interfaces.Group, userIDs []int64, otx ...*sql.Tx) error{
	panic("unimplemented")
}
func (*MockGroupModel) DeleteGroup(ctx context.Context, groupID int64, requestingUserID int64, otx ...*sql.Tx) error{
	panic("unimplemented")
}
func (*MockGroupModel) GetGroupByID(ctx context.Context, groupID int64, requestingUserID int64, otx ...*sql.Tx) (*interfaces.Group, error){
	panic("unimplemented")
}
func (*MockGroupModel) GetGroupByScope(ctx context.Context, scopeID int64, requestingUserID int64, otx ...*sql.Tx) (*interfaces.Group, error){
	panic("unimplemented")
}
var _ interfaces.GroupService = &MockGroupModel{}

// Implement methods like InsertGroup, UpdateGroup, etc., similar to the MockCategoryModel methods.
