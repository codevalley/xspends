package impl

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type UserScopeModel struct {
	TableUserScopes string
	ColumnUserID    string
	ColumnScopeID   string
	ColumnRole      string
}

func NewUserScopeModel() *UserScopeModel {
	return &UserScopeModel{
		TableUserScopes: "user_scopes",
		ColumnUserID:    "user_id",
		ColumnScopeID:   "scope_id",
		ColumnRole:      "role",
	}
}

// UpsertUserScope either inserts a new user-scope relationship or updates an existing one.
func (usm *UserScopeModel) UpsertUserScope(ctx context.Context, userID, scopeID int64, role string, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	query, args, err := GetQueryBuilder().
		Insert(usm.TableUserScopes).
		Columns(usm.ColumnUserID, usm.ColumnScopeID, usm.ColumnRole).
		Values(userID, scopeID, role).
		Suffix("ON DUPLICATE KEY UPDATE " + usm.ColumnRole + " = VALUES(" + usm.ColumnRole + ")").
		ToSql()
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "building upsert query failed")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "executing upsert query failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}

// GetUserScope retrieves a specific user-scope relationship.
func (usm *UserScopeModel) GetUserScope(ctx context.Context, userID, scopeID int64, otx ...*sql.Tx) (*interfaces.UserScope, error) {
	_, executor := getExecutor(otx...)

	query, args, err := GetQueryBuilder().
		Select(usm.ColumnUserID, usm.ColumnScopeID, usm.ColumnRole).
		From(usm.TableUserScopes).
		Where(squirrel.Eq{usm.ColumnUserID: userID, usm.ColumnScopeID: scopeID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building select query failed")
	}

	var userScope interfaces.UserScope
	err = executor.QueryRowContext(ctx, query, args...).Scan(&userScope.UserID, &userScope.ScopeID, &userScope.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user-scope relationship not found")
		}
		return nil, errors.Wrap(err, "querying user-scope relationship failed")
	}

	return &userScope, nil
}

// GetUserScope retrieves a specific user-scope relationship.
func (usm *UserScopeModel) ValidateUserScope(ctx context.Context, userID, scopeID int64, role string, otx ...*sql.Tx) bool {
	_, executor := getExecutor(otx...)

	// Determine the minimum access level required based on the provided role
	minAccessLevel, ok := roleHierarchy[role]
	if !ok {
		return false
	}

	// Build a slice of roles to include based on the hierarchy
	var rolesToInclude []string
	for r, level := range roleHierarchy {
		if level <= minAccessLevel {
			rolesToInclude = append(rolesToInclude, r)
		}
	}

	query, args, err := GetQueryBuilder().
		Select(usm.ColumnUserID, usm.ColumnScopeID, usm.ColumnRole).
		From(usm.TableUserScopes).
		Where(squirrel.Eq{usm.ColumnUserID: userID, usm.ColumnScopeID: scopeID}).
		Where(squirrel.Eq{usm.ColumnRole: rolesToInclude}).
		ToSql()
	if err != nil {
		return false
	}

	var userScope interfaces.UserScope
	err = executor.QueryRowContext(ctx, query, args...).Scan(&userScope.UserID, &userScope.ScopeID, &userScope.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}

	return true
}

// GetUserScope retrieves user-scope relationships based on access level hierarchy.
func (usm *UserScopeModel) GetUserScopesByRole(ctx context.Context, userID int64, role string, otx ...*sql.Tx) ([]interfaces.UserScope, error) {
	_, executor := getExecutor(otx...)

	// Determine the minimum access level required based on the provided role
	minAccessLevel, ok := roleHierarchy[role]
	if !ok {
		return nil, errors.Errorf("invalid role: %s", role)
	}

	// Build a slice of roles to include based on the hierarchy
	var rolesToInclude []string
	for r, level := range roleHierarchy {
		if level <= minAccessLevel {
			rolesToInclude = append(rolesToInclude, r)
		}
	}

	// Construct the query to select user scopes where the role is in the list of roles to include
	query, args, err := GetQueryBuilder().
		Select(usm.ColumnUserID, usm.ColumnScopeID, usm.ColumnRole).
		From(usm.TableUserScopes).
		Where(squirrel.Eq{usm.ColumnUserID: userID}).
		Where(squirrel.Eq{usm.ColumnRole: rolesToInclude}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building select query failed")
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying user-scope relationships failed")
	}
	defer rows.Close()

	var userScopes []interfaces.UserScope
	for rows.Next() {
		var userScope interfaces.UserScope
		if err := rows.Scan(&userScope.UserID, &userScope.ScopeID, &userScope.Role); err != nil {
			return nil, errors.Wrap(err, "scanning user-scope row failed")
		}
		userScopes = append(userScopes, userScope)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "processing user-scope rows failed")
	}

	return userScopes, nil
}

// DeleteUserScope removes a user-scope relationship.
func (usm *UserScopeModel) DeleteUserScope(ctx context.Context, userID, scopeID int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	query, args, err := GetQueryBuilder().
		Delete(usm.TableUserScopes).
		Where(squirrel.Eq{usm.ColumnUserID: userID, usm.ColumnScopeID: scopeID}).
		ToSql()
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "building delete query failed")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "executing delete query failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}
