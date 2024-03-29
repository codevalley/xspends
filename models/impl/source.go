/*
MIT License

Copyright (c) 2023 Narayan Babu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package impl

import (
	"context"
	"database/sql"
	"strings"
	"time"
	"xspends/models/interfaces"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type SourceModel struct {
	TableSources      string
	ColumnID          string
	ColumnUserID      string
	ColumnName        string
	ColumnType        string
	ColumnBalance     string
	ColumnScope       string
	ColumnCreatedAt   string
	ColumnUpdatedAt   string
	SourceTypeCredit  string
	SourceTypeSavings string
}

func NewSourceModel() *SourceModel {
	return &SourceModel{
		TableSources:      "sources",
		ColumnID:          "source_id",
		ColumnUserID:      "user_id",
		ColumnName:        "name",
		ColumnType:        "type",
		ColumnBalance:     "balance",
		ColumnScope:       "scope_id",
		ColumnCreatedAt:   "created_at",
		ColumnUpdatedAt:   "updated_at",
		SourceTypeCredit:  "CREDIT",
		SourceTypeSavings: "SAVINGS",
	}
}

func (cm *SourceModel) validateSourceInput(ctx context.Context, source *interfaces.Source, role string) error {
	if source.ScopeID <= 0 || source.UserID <= 0 || source.Name == "" {
		return errors.New(ErrInvalidInput)
	}

	if !strings.EqualFold(source.Type, cm.SourceTypeCredit) && !strings.EqualFold(source.Type, cm.SourceTypeSavings) {
		return errors.New("invalid type: type must be CREDIT or SAVINGS")
	}

	if !GetModelsService().UserScopeModel.ValidateUserScope(ctx, source.UserID, source.ScopeID, role) {
		return errors.New(ErrInvalidInput)
	}
	return nil
}

func (sm *SourceModel) InsertSource(ctx context.Context, source *interfaces.Source, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	if err := sm.validateSourceInput(ctx, source, RoleWrite); err != nil {
		return err
	}

	var err error
	source.ID, err = util.GenerateSnowflakeID() //no checks here. If it fails, it fails.
	if err != nil {
		return errors.Wrap(err, "generating Snowflake ID failed")
	}
	source.CreatedAt = time.Now()
	source.UpdatedAt = source.CreatedAt

	query, args, err := GetQueryBuilder().Insert(sm.TableSources).
		Columns(sm.ColumnID, sm.ColumnUserID, sm.ColumnName, sm.ColumnType, sm.ColumnBalance, sm.ColumnScope, sm.ColumnCreatedAt, sm.ColumnUpdatedAt).
		Values(source.ID, source.UserID, source.Name, source.Type, source.Balance, source.ScopeID, source.CreatedAt, source.UpdatedAt).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "preparing insert SQL for source")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing insert for source")
	}
	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (sm *SourceModel) UpdateSource(ctx context.Context, source *interfaces.Source, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	if err := sm.validateSourceInput(ctx, source, RoleWrite); err != nil {
		return err
	}

	source.UpdatedAt = time.Now()

	query, args, err := GetQueryBuilder().Update(sm.TableSources).
		Set(sm.ColumnName, source.Name).
		Set(sm.ColumnType, source.Type).
		Set(sm.ColumnBalance, source.Balance).
		Set(sm.ColumnUpdatedAt, source.UpdatedAt).
		Where(squirrel.Eq{sm.ColumnID: source.ID, sm.ColumnUserID: source.ScopeID}).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "preparing update SQL for source")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing update for source")
	}
	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (sm *SourceModel) DeleteSource(ctx context.Context, sourceID int64, scopes []int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)
	query, args, err := GetQueryBuilder().Delete(sm.TableSources).
		Where(squirrel.Eq{sm.ColumnID: sourceID, sm.ColumnScope: scopes}).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "preparing delete SQL for source")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing delete for source")
	}
	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (sm *SourceModel) GetSourceByID(ctx context.Context, sourceID int64, scopes []int64, otx ...*sql.Tx) (*interfaces.Source, error) {
	_, executor := getExecutor(otx...)
	query, args, err := GetQueryBuilder().Select(sm.ColumnID, sm.ColumnUserID, sm.ColumnName, sm.ColumnType, sm.ColumnBalance, sm.ColumnScope, sm.ColumnCreatedAt, sm.ColumnUpdatedAt).
		From(sm.TableSources).
		Where(squirrel.Eq{sm.ColumnID: sourceID, sm.ColumnScope: scopes}).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "preparing select SQL for source by ID")
	}

	source := &interfaces.Source{}
	err = executor.QueryRowContext(ctx, query, args...).Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance, &source.ScopeID, &source.CreatedAt, &source.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("source not found")
		}
		return nil, errors.Wrap(err, "querying source by ID")
	}

	return source, nil
}

// paginated
func (sm *SourceModel) GetScopedSources(ctx context.Context, page int, itemsPerPage int, scopes []int64, otx ...*sql.Tx) ([]interfaces.Source, error) {
	_, executor := getExecutor(otx...)

	offset := (page - 1) * itemsPerPage

	query, args, err := GetQueryBuilder().Select(sm.ColumnID, sm.ColumnUserID, sm.ColumnName, sm.ColumnType, sm.ColumnBalance, sm.ColumnScope, sm.ColumnCreatedAt, sm.ColumnUpdatedAt).
		From(sm.TableSources).
		Where(squirrel.Eq{sm.ColumnScope: scopes}).
		Limit(uint64(itemsPerPage)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "preparing paginated select statement for sources failed")
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying paginated sources failed")
	}
	defer rows.Close()

	var sources []interfaces.Source
	for rows.Next() {
		var source interfaces.Source
		if err = rows.Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance, &source.ScopeID, &source.CreatedAt, &source.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "scanning paginated source row failed")
		}
		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "during row processing for sources")
	}

	return sources, nil
}

func (sm *SourceModel) GetSources(ctx context.Context, scopes []int64, otx ...*sql.Tx) ([]interfaces.Source, error) {
	_, executor := getExecutor(otx...)
	query, args, err := GetQueryBuilder().Select(sm.ColumnID, sm.ColumnUserID, sm.ColumnName, sm.ColumnType, sm.ColumnBalance, sm.ColumnScope, sm.ColumnCreatedAt, sm.ColumnUpdatedAt).
		From(sm.TableSources).
		Where(squirrel.Eq{sm.ColumnScope: scopes}).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "preparing select SQL for sources by user ID")
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying sources by user ID")
	}
	defer rows.Close()

	var sources []interfaces.Source
	for rows.Next() {
		var source interfaces.Source
		if err = rows.Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance, &source.ScopeID, &source.CreatedAt, &source.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "scanning source row")
		}
		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "during row processing for sources")
	}

	return sources, nil
}

func (sm *SourceModel) SourceIDExists(ctx context.Context, sourceID int64, scopes []int64, otx ...*sql.Tx) (bool, error) {
	_, executor := getExecutor(otx...)

	query, args, err := GetQueryBuilder().Select("1").
		From(sm.TableSources).
		Where(squirrel.Eq{sm.ColumnID: sourceID, sm.ColumnScope: scopes}).
		Limit(1).
		ToSql()

	if err != nil {
		return false, errors.Wrap(err, "preparing SQL to check if source exists by ID")
	}

	var exists int
	err = executor.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "checking if source exists by ID")
	}

	return exists == 1, nil
}
