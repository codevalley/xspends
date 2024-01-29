package impl

import (
	"context"
	"database/sql"
	"time"
	"xspends/models/interfaces"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type GroupModel struct {
	TableGroups               string
	ColumnGroupID             string
	ColumnOwnerID             string
	ColumnScopeID             string
	ColumnGroupName           string
	ColumnDescription         string
	ColumnIcon                string
	ColumnStatus              string
	ColumnCreatedAt           string
	ColumnUpdatedAt           string
	MaxGroupNameLength        int
	MaxGroupDescriptionLength int
}

func NewGroupModel() *GroupModel {
	return &GroupModel{
		TableGroups:               "groups",
		ColumnGroupID:             "group_id",
		ColumnOwnerID:             "owner_id",
		ColumnScopeID:             "scope_id",
		ColumnGroupName:           "group_name",
		ColumnDescription:         "description",
		ColumnIcon:                "icon",
		ColumnStatus:              "status",
		ColumnCreatedAt:           "created_at",
		ColumnUpdatedAt:           "updated_at",
		MaxGroupNameLength:        100, // Adjust as per your requirement
		MaxGroupDescriptionLength: 512, // Adjust as per your requirement
	}
}
func (gm *GroupModel) validateGroupInput(group *interfaces.Group) error {
	if group.OwnerID <= 0 || group.GroupName == "" || len(group.GroupName) > gm.MaxGroupNameLength || len(group.Description) > gm.MaxGroupDescriptionLength {
		return errors.New(ErrInvalidInput)
	}
	return nil
}
func (gm *GroupModel) CreateGroup(ctx context.Context, group *interfaces.Group, userIDs []int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	if err := gm.validateGroupInput(group); err != nil {
		return err
	}

	// Initialize ScopeModel and create a new scope
	scopeID, err := GetModelsService().ScopeModel.CreateScope(ctx, "group", otx...)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "creating new scope failed")
	}

	// Generate Snowflake ID for the group
	group.GroupID, err = util.GenerateSnowflakeID()
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "generating Snowflake ID for group failed")
	}

	group.CreatedAt, group.UpdatedAt = time.Now(), time.Now()

	// Insert into groups table
	groupsQuery, groupsArgs, err := GetQueryBuilder().Insert(gm.TableGroups).
		Columns(gm.ColumnGroupID, gm.ColumnOwnerID, gm.ColumnScopeID, gm.ColumnGroupName, gm.ColumnDescription, gm.ColumnIcon, gm.ColumnStatus, gm.ColumnCreatedAt, gm.ColumnUpdatedAt).
		Values(group.GroupID, group.OwnerID, scopeID, group.GroupName, group.Description, group.Icon, group.Status, group.CreatedAt, group.UpdatedAt).
		ToSql()
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "building groups insert query failed")
	}

	_, err = executor.ExecContext(ctx, groupsQuery, groupsArgs...)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "inserting into groups failed")
	}

	//TODO: To be separated out to scope insert
	// Link users to the group's scope
	for _, userID := range userIDs {
		userScopesQuery, userScopesArgs, err := GetQueryBuilder().Insert("user_scopes").
			Columns("user_id", "scope_id").
			Values(userID, scopeID).
			ToSql()
		if err != nil {
			commitOrRollback(executor, isExternalTx, err)
			return errors.Wrap(err, "building user_scopes insert query failed")
		}

		_, err = executor.ExecContext(ctx, userScopesQuery, userScopesArgs...)
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

	groupDeleteQuery, groupDeleteArgs, err := GetQueryBuilder().Delete(gm.TableGroups).
		Where(squirrel.Eq{gm.ColumnGroupID: groupID}).
		ToSql()
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "building group delete query failed")
	}

	_, err = executor.ExecContext(ctx, groupDeleteQuery, groupDeleteArgs...)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "deleting group failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (gm *GroupModel) GetGroupByID(ctx context.Context, groupID int64, requestingUserID int64, otx ...*sql.Tx) (*interfaces.Group, error) {
	_, executor := getExecutor(otx...)

	// Fetch group details
	groupSelectQuery, groupSelectArgs, err := GetQueryBuilder().Select(gm.ColumnGroupID, gm.ColumnOwnerID, gm.ColumnScopeID, gm.ColumnGroupName, gm.ColumnDescription, gm.ColumnIcon, gm.ColumnStatus, gm.ColumnCreatedAt, gm.ColumnUpdatedAt).
		From(gm.TableGroups).
		Where(squirrel.Eq{gm.ColumnGroupID: groupID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building group select query failed")
	}

	row := executor.QueryRowContext(ctx, groupSelectQuery, groupSelectArgs...)
	group := interfaces.Group{}
	if err := row.Scan(&group.GroupID, &group.OwnerID, &group.ScopeID, &group.GroupName, &group.Description, &group.Icon, &group.Status, &group.CreatedAt, &group.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("group not found")
		}
		return nil, errors.Wrap(err, "querying group by ID failed")
	}

	return &group, nil
}

func (gm *GroupModel) GetGroupByScope(ctx context.Context, scopeID int64, requestingUserID int64, otx ...*sql.Tx) (*interfaces.Group, error) {
	_, executor := getExecutor(otx...)

	// Fetch group details by scope ID
	groupSelectQuery, groupSelectArgs, err := GetQueryBuilder().Select(gm.ColumnGroupID, gm.ColumnOwnerID, gm.ColumnScopeID, gm.ColumnGroupName, gm.ColumnDescription, gm.ColumnIcon, gm.ColumnStatus, gm.ColumnCreatedAt, gm.ColumnUpdatedAt).
		From(gm.TableGroups).
		Where(squirrel.Eq{gm.ColumnScopeID: scopeID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building group select by scope query failed")
	}

	row := executor.QueryRowContext(ctx, groupSelectQuery, groupSelectArgs...)
	group := interfaces.Group{}
	if err := row.Scan(&group.GroupID, &group.OwnerID, &group.ScopeID, &group.GroupName, &group.Description, &group.Icon, &group.Status, &group.CreatedAt, &group.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("group not found for the given scope")
		}
		return nil, errors.Wrap(err, "querying group by scope failed")
	}

	return &group, nil
}

//Add to groups method
//get users in a group method
