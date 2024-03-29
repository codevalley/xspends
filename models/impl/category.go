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
	"time"
	"xspends/models/interfaces"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

const ErrInvalidInput = "invalid input: user ID must be numeric, name must not be empty or exceed max length, description must not exceed max length"
const ErrInvalidScope = "Invalid scope presented for the request"

type CategoryModel struct {
	TableCategories              string
	ColumnID                     string
	ColumnUserID                 string
	ColumnName                   string
	ColumnDescription            string
	ColumnIcon                   string
	ColumnScopeID                string
	ColumnCreatedAt              string
	ColumnUpdatedAt              string
	MaxCategoryNameLength        int
	MaxCategoryDescriptionLength int
}

func NewCategoryModel() *CategoryModel {
	return &CategoryModel{
		TableCategories:              "categories",
		ColumnID:                     "category_id",
		ColumnUserID:                 "user_id",
		ColumnName:                   "name",
		ColumnDescription:            "description",
		ColumnIcon:                   "icon",
		ColumnScopeID:                "scope_id",
		ColumnCreatedAt:              "created_at",
		ColumnUpdatedAt:              "updated_at",
		MaxCategoryNameLength:        100, // Adjust as per your requirement
		MaxCategoryDescriptionLength: 512, // Adjust as per your requirement
	}
}

func (cm *CategoryModel) validateCategoryInput(ctx context.Context, category *interfaces.Category, role string) error {
	if category.ScopeID <= 0 || category.UserID <= 0 || category.Name == "" || len(category.Name) > cm.MaxCategoryNameLength || len(category.Description) > cm.MaxCategoryDescriptionLength {
		return errors.New(ErrInvalidInput)
	}

	if !GetModelsService().UserScopeModel.ValidateUserScope(ctx, category.UserID, category.ScopeID, role) {
		return errors.New(ErrInvalidScope)
	}
	return nil
}

// InsertCategory inserts a new category into the database.
func (cm *CategoryModel) InsertCategory(ctx context.Context, category *interfaces.Category, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	if err := cm.validateCategoryInput(ctx, category, RoleWrite); err != nil {
		return err
	}

	var err error
	category.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		return errors.Wrap(err, "generating Snowflake ID failed")
	}
	category.CreatedAt, category.UpdatedAt = time.Now(), time.Now()

	query, args, err := GetQueryBuilder().Insert(cm.TableCategories).
		Columns(cm.ColumnID, cm.ColumnUserID, cm.ColumnName, cm.ColumnDescription, cm.ColumnIcon, cm.ColumnScopeID, cm.ColumnCreatedAt, cm.ColumnUpdatedAt).
		Values(category.ID, category.UserID, category.Name, category.Description, category.Icon, category.ScopeID, category.CreatedAt, category.UpdatedAt).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "preparing insert statement failed")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing insert statement failed")
	}

	commitOrRollback(executor, isExternalTx, err)

	return nil
}

// UpdateCategory updates an existing category in the database.
func (cm *CategoryModel) UpdateCategory(ctx context.Context, category *interfaces.Category, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	if err := cm.validateCategoryInput(ctx, category, RoleWrite); err != nil {
		return err
	}

	category.UpdatedAt = time.Now()

	query, args, err := GetQueryBuilder().Update(cm.TableCategories).
		Set(cm.ColumnName, category.Name).
		Set(cm.ColumnDescription, category.Description).
		Set(cm.ColumnIcon, category.Icon).
		Set(cm.ColumnUpdatedAt, category.UpdatedAt).
		Where(squirrel.Eq{cm.ColumnID: category.ID, cm.ColumnScopeID: category.ScopeID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "preparing update statement failed")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing update statement failed")
	}
	commitOrRollback(executor, isExternalTx, err)

	return nil
}

// DeleteCategory deletes a category from the database.
func (cm *CategoryModel) DeleteCategory(ctx context.Context, categoryID int64, scopes []int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	// if !GetModelsService().UserScopeModel.ValidateUserScope(ctx, category.UserID, scopes, RoleWrite) {
	// 	return errors.New(ErrInvalidScope)
	// }
	//TODO: No check if the current user has access to the scope (will happen in handler, but no double check)
	//We can add validation here as well
	query, args, err := GetQueryBuilder().Delete(cm.TableCategories).
		Where(squirrel.Eq{cm.ColumnID: categoryID, cm.ColumnScopeID: scopes}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "preparing delete statement failed")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing delete statement failed")
	}
	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (cm *CategoryModel) GetCategoryByID(ctx context.Context, categoryID int64, scopes []int64, otx ...*sql.Tx) (*interfaces.Category, error) {
	_, executor := getExecutor(otx...)

	query, args, err := sqlBuilder.Select(cm.ColumnID, cm.ColumnUserID, cm.ColumnScopeID, cm.ColumnName, cm.ColumnDescription, cm.ColumnIcon, cm.ColumnScopeID, cm.ColumnCreatedAt, cm.ColumnUpdatedAt).
		From(cm.TableCategories).
		Where(squirrel.Eq{cm.ColumnID: categoryID, cm.ColumnScopeID: scopes}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "preparing select statement for a category by ID failed")
	}

	var category interfaces.Category
	err = executor.QueryRowContext(ctx, query, args...).Scan(&category.ID, &category.UserID, &category.ScopeID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, errors.Wrap(err, "querying category by ID failed")
	}

	return &category, nil
}

func (cm *CategoryModel) GetScopedCategories(ctx context.Context, page int, itemsPerPage int, scopes []int64, otx ...*sql.Tx) ([]interfaces.Category, error) {
	_, executor := getExecutor(otx...)

	offset := (page - 1) * itemsPerPage

	query, args, err := sqlBuilder.Select(cm.ColumnID, cm.ColumnUserID, cm.ColumnName, cm.ColumnDescription, cm.ColumnIcon, cm.ColumnScopeID, cm.ColumnCreatedAt, cm.ColumnUpdatedAt).
		From(cm.TableCategories).
		Where(squirrel.Eq{cm.ColumnUserID: scopes}).
		Limit(uint64(itemsPerPage)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "preparing paginated select statement for categories failed")
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying paginated categories failed")
	}
	defer rows.Close()

	var categories []interfaces.Category
	for rows.Next() {
		category := interfaces.Category{}
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "scanning paginated category row failed")
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (cm *CategoryModel) CategoryIDExists(ctx context.Context, categoryID int64, scopes []int64, otx ...*sql.Tx) (bool, error) {
	_, executor := getExecutor(otx...)
	query, args, err := sqlBuilder.Select("1").
		From(cm.TableCategories).
		Where(squirrel.Eq{cm.ColumnID: categoryID, cm.ColumnScopeID: scopes}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, errors.Wrap(err, "preparing statement to check category existence failed")
	}

	var exists int
	err = executor.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "checking category existence failed")
	}
	return exists == 1, nil
}
