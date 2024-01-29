package impl

import (
	"context"
	"database/sql"
	"xspends/models/interfaces"
	"xspends/util"

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

	_, err = executor.ExecContext(ctx, "INSERT INTO scopes (scope_id, type) VALUES (?, ?)", scopeID, scopeType)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return 0, errors.Wrap(err, "inserting into scopes failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return scopeID, nil
}

func (sm *ScopeModel) GetScope(ctx context.Context, scopeID int64, otx ...*sql.Tx) (*interfaces.Scope, error) {
	_, executor := getExecutor(otx...)

	row := executor.QueryRowContext(ctx, "SELECT scope_id, type FROM scopes WHERE scope_id = ?", scopeID)
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

	_, err := executor.ExecContext(ctx, "DELETE FROM scopes WHERE scope_id = ?", scopeID)
	if err != nil {
		commitOrRollback(executor, isExternalTx, err)
		return errors.Wrap(err, "deleting scope failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}
