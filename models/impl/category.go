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

const ErrInvalidInput = "invalid input: user ID must be positive, name must not be empty or exceed max length, description must not exceed max length"

type CategoryModel struct {
	TableCategories              string
	ColumnID                     string
	ColumnUserID                 string
	ColumnName                   string
	ColumnDescription            string
	ColumnIcon                   string
	ColumnCreatedAt              string
	ColumnUpdatedAt              string
	MaxCategoryNameLength        int
	MaxCategoryDescriptionLength int
}

func NewCategoryModel() *CategoryModel {
	return &CategoryModel{
		TableCategories:              "categories",
		ColumnID:                     "id",
		ColumnUserID:                 "user_id",
		ColumnName:                   "name",
		ColumnDescription:            "description",
		ColumnIcon:                   "icon",
		ColumnCreatedAt:              "created_at",
		ColumnUpdatedAt:              "updated_at",
		MaxCategoryNameLength:        100, // Adjust as per your requirement
		MaxCategoryDescriptionLength: 512, // Adjust as per your requirement
	}
}

func (cm *CategoryModel) validateCategoryInput(category *interfaces.Category) error {
	if category.UserID <= 0 || category.Name == "" || len(category.Name) > cm.MaxCategoryNameLength || len(category.Description) > cm.MaxCategoryDescriptionLength {
		return errors.New(ErrInvalidInput)
	}
	return nil
}

// InsertCategory inserts a new category into the database.
func (cm *CategoryModel) InsertCategory(ctx context.Context, category *interfaces.Category, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	if err := cm.validateCategoryInput(category); err != nil {
		return err
	}

	var err error
	category.ID, _ = util.GenerateSnowflakeID()
	if err != nil {
		return errors.Wrap(err, "generating Snowflake ID failed")
	}
	category.CreatedAt, category.UpdatedAt = time.Now(), time.Now()

	query, args, err := GetQueryBuilder().Insert("categories").
		Columns("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		Values(category.ID, category.UserID, category.Name, category.Description, category.Icon, category.CreatedAt, category.UpdatedAt).
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

	if err := cm.validateCategoryInput(category); err != nil {
		return err
	}

	category.UpdatedAt = time.Now()

	query, args, err := GetQueryBuilder().Update("categories").
		Set("name", category.Name).
		Set("description", category.Description).
		Set("icon", category.Icon).
		Set("updated_at", category.UpdatedAt).
		Where(squirrel.Eq{"id": category.ID, "user_id": category.UserID}).
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
func (cm *CategoryModel) DeleteCategory(ctx context.Context, categoryID int64, userID int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	query, args, err := GetQueryBuilder().Delete("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
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

// GetAllCategories retrieves all categories for a user from the database.
func (cm *CategoryModel) GetAllCategories(ctx context.Context, userID int64, otx ...*sql.Tx) ([]interfaces.Category, error) {
	_, executor := getExecutor(otx...)

	query, args, err := GetQueryBuilder().Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "preparing select statement for all categories failed")
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying categories failed")
	}
	defer rows.Close()

	var categories []interfaces.Category
	for rows.Next() {
		category := interfaces.Category{}
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "scanning category row failed")
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by its ID for a user from the database.
func (cm *CategoryModel) GetCategoryByID(ctx context.Context, categoryID int64, userID int64, otx ...*sql.Tx) (*interfaces.Category, error) {
	_, executor := getExecutor(otx...)

	query, args, err := sqlBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "preparing select statement for a category by ID failed")
	}

	var category interfaces.Category
	err = executor.QueryRowContext(ctx, query, args...).Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, errors.Wrap(err, "querying category by ID failed")
	}

	return &category, nil
}

// GetPagedCategories retrieves a paginated list of categories for a user from the database.
func (cm *CategoryModel) GetPagedCategories(ctx context.Context, page int, itemsPerPage int, userID int64, otx ...*sql.Tx) ([]interfaces.Category, error) {
	_, executor := getExecutor(otx...)

	offset := (page - 1) * itemsPerPage

	query, args, err := sqlBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"user_id": userID}).
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

// CategoryIDExists checks if a category with the given ID exists in the database.
func (cm *CategoryModel) CategoryIDExists(ctx context.Context, categoryID int64, userID int64, otx ...*sql.Tx) (bool, error) {
	_, executor := getExecutor(otx...)
	query, args, err := sqlBuilder.Select("1").
		From("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
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
