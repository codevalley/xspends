package impl

import (
	"context"
	"database/sql"
	"time"
	"xspends/models/interfaces"
	"xspends/util"

	"github.com/pkg/errors"
)

type GroupModel struct {
	TableGroups       string
	ColumnGroupID     string
	ColumnOwnerID     string
	ColumnScopeID     string
	ColumnGroupName   string
	ColumnDescription string
	ColumnIcon        string
	ColumnStatus      string
	ColumnCreatedAt   string
	ColumnUpdatedAt   string
}

func NewGroupModel() *GroupModel {
	return &GroupModel{
		TableGroups:       "groups",
		ColumnGroupID:     "group_id",
		ColumnOwnerID:     "owner_id",
		ColumnScopeID:     "scope_id",
		ColumnGroupName:   "group_name",
		ColumnDescription: "description",
		ColumnIcon:        "icon",
		ColumnStatus:      "status",
		ColumnCreatedAt:   "created_at",
		ColumnUpdatedAt:   "updated_at",
	}
}

func (gm *GroupModel) CreateGroup(ctx context.Context, group *interfaces.Group, userIDs []int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	scopeID, _ := util.GenerateSnowflakeID()      // Error handling for ID generation is required
	group.GroupID, _ = util.GenerateSnowflakeID() // Error handling for ID generation is required
	group.CreatedAt, group.UpdatedAt = time.Now(), time.Now()

	// Insert into scopes table
	_, err := executor.ExecContext(ctx, "INSERT INTO scopes (scope_id, type) VALUES (?, ?)", scopeID, "group")
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "inserting into scopes failed")
	}

	// Insert into groups table
	_, err = executor.ExecContext(ctx, "INSERT INTO groups (group_id, owner_id, scope_id, group_name, description, icon, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		group.GroupID, group.OwnerID, scopeID, group.GroupName, group.Description, group.Icon, group.Status, group.CreatedAt, group.UpdatedAt)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "inserting into groups failed")
	}

	// Link users to the group's scope
	for _, userID := range userIDs {
		_, err = executor.ExecContext(ctx, "INSERT INTO user_scopes (user_id, scope_id) VALUES (?, ?)", userID, scopeID)
		if err != nil {
			commitOrRollback(executor, isExternalTx, err)
			return errors.Wrap(err, "inserting into user_scopes failed")
		}
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}
