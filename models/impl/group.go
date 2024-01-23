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

func (gm *GroupModel) DeleteGroup(ctx context.Context, groupID int64, requestingUserID int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	// Verify ownership
	row := executor.QueryRowContext(ctx, "SELECT owner_id FROM groups WHERE group_id = ?", groupID)
	var ownerID int64
	if err := row.Scan(&ownerID); err != nil {
		commitOrRollback(executor, isExternalTx, err)
		if err == sql.ErrNoRows {
			return errors.New("group not found")
		}
		return errors.Wrap(err, "verifying group ownership failed")
	}
	if ownerID != requestingUserID {
		return errors.New("unauthorized to delete group")
	}

	// Delete group and associated scope
	_, err := executor.ExecContext(ctx, "DELETE FROM groups WHERE group_id = ?", groupID)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "deleting group failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (gm *GroupModel) GetGroupByID(ctx context.Context, groupID int64, requestingUserID int64, otx ...*sql.Tx) (*interfaces.Group, error) {
	_, executor := getExecutor(otx...)

	// Ensure user has access
	row := executor.QueryRowContext(ctx, "SELECT 1 FROM user_scopes WHERE user_id = ? AND scope_id = (SELECT scope_id FROM groups WHERE group_id = ?)", requestingUserID, groupID)
	var exists int
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("access to group not found or group does not exist")
		}
		return nil, errors.Wrap(err, "verifying access to group failed")
	}

	// Fetch group details
	row = executor.QueryRowContext(ctx, "SELECT group_id, owner_id, scope_id, group_name, description, icon, status, created_at, updated_at FROM groups WHERE group_id = ?", groupID)
	group := interfaces.Group{}
	if err := row.Scan(&group.GroupID, &group.OwnerID, &group.ScopeID, &group.GroupName, &group.Description, &group.Icon, &group.Status, &group.CreatedAt, &group.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("group not found")
		}
		return nil, errors.Wrap(err, "querying group by ID failed")
	}

	return &group, nil
}
