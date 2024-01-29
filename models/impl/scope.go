package impl

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type ScopeModel struct {
	TableScopes   string
	ColumnScopeID string
	ColumnType    string
}

func NewScopeModel() *ScopeModel {
	return &ScopeModel{
		TableScopes:   "scopes",
		ColumnScopeID: "scope_id",
		ColumnType:    "type",
	}
}

func (sm *ScopeModel) CreateScope(ctx context.Context, scopeType string, otx ...*sql.Tx) (int64, error) {
	isExternalTx, executor := getExecutor(otx...)

	scopeID, err := util.GenerateSnowflakeID()
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return 0, errors.Wrap(err, "generating Snowflake ID failed")
	}

	insertQuery, args, err := GetQueryBuilder().Insert(sm.TableScopes).
		Columns(sm.ColumnScopeID, sm.ColumnType).
		Values(scopeID, scopeType).
		ToSql()
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return 0, errors.Wrap(err, "building insert query failed")
	}

	_, err = executor.ExecContext(ctx, insertQuery, args...)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return 0, errors.Wrap(err, "inserting into scopes failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return scopeID, nil
}

func (sm *ScopeModel) GetScope(ctx context.Context, scopeID int64, otx ...*sql.Tx) (*interfaces.Scope, error) {
	_, executor := getExecutor(otx...)

	selectQuery, args, err := GetQueryBuilder().Select(sm.ColumnScopeID, sm.ColumnType).
		From(sm.TableScopes).
		Where(squirrel.Eq{sm.ColumnScopeID: scopeID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building select query failed")
	}

	row := executor.QueryRowContext(ctx, selectQuery, args...)
	scope := &interfaces.Scope{}
	if err := row.Scan(&scope.ScopeID, &scope.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("scope not found")
		}
		return nil, errors.Wrap(err, "querying scope failed")
	}

	return scope, nil
}

func (sm *ScopeModel) DeleteScope(ctx context.Context, scopeID int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	deleteQuery, args, err := GetQueryBuilder().Delete(sm.TableScopes).
		Where(squirrel.Eq{sm.ColumnScopeID: scopeID}).
		ToSql()
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "building delete query failed")
	}

	_, err = executor.ExecContext(ctx, deleteQuery, args...)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "deleting scope failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}
