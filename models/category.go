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

package models

import (
	"context"
	"database/sql"
	"time"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

const (
	maxCategoryNameLength        = 255
	maxCategoryDescriptionLength = 500
	ErrInvalidInput              = "invalid input: user ID must be positive, name must not be empty or exceed max length, description must not exceed max length"
)

type Category struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func validateCategoryInput(category *Category) error {
	if category.UserID <= 0 || category.Name == "" || len(category.Name) > maxCategoryNameLength || len(category.Description) > maxCategoryDescriptionLength {
		return errors.New(ErrInvalidInput)
	}
	return nil
}

// InsertCategory inserts a new category into the database.
func InsertCategory(ctx context.Context, category *Category, dbService *DBService, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(dbService, otx...)

	if err := validateCategoryInput(category); err != nil {
		return err
	}

	category.ID, _ = util.GenerateSnowflakeID()
	category.CreatedAt, category.UpdatedAt = time.Now(), time.Now()

	query, args, err := SQLBuilder.Insert("categories").
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

	if !isExternalTx {
		if tx, ok := executor.(*sql.Tx); ok {
			if err := tx.Commit(); err != nil {
				tx.Rollback()
				return errors.Wrap(err, "committing transaction failed")
			}
		}
	}

	return nil
}

// UpdateCategory updates an existing category in the database.
func UpdateCategory(ctx context.Context, category *Category, dbService *DBService, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(dbService, otx...)

	if err := validateCategoryInput(category); err != nil {
		return err
	}

	category.UpdatedAt = time.Now()

	query, args, err := SQLBuilder.Update("categories").
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
	if !isExternalTx {
		if tx, ok := executor.(*sql.Tx); ok {
			if err := tx.Commit(); err != nil {
				tx.Rollback()
				return errors.Wrap(err, "committing transaction failed")
			}
		}
	}

	return nil
}

// DeleteCategory deletes a category from the database.
func DeleteCategory(ctx context.Context, categoryID int64, userID int64, dbService *DBService, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(dbService, otx...)

	query, args, err := SQLBuilder.Delete("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "preparing delete statement failed")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing delete statement failed")
	}
	if !isExternalTx {
		if tx, ok := executor.(*sql.Tx); ok {
			if err := tx.Commit(); err != nil {
				tx.Rollback()
				return errors.Wrap(err, "committing transaction failed")
			}
		}
	}
	return nil
}

// GetAllCategories retrieves all categories for a user from the database.
func GetAllCategories(ctx context.Context, userID int64, dbService *DBService, otx ...*sql.Tx) ([]Category, error) {
	_, executor := getExecutor(dbService, otx...)

	query, args, err := SQLBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
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

	var categories []Category
	for rows.Next() {
		category := Category{}
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "scanning category row failed")
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by its ID for a user from the database.
func GetCategoryByID(ctx context.Context, categoryID int64, userID int64, dbService *DBService, otx ...*sql.Tx) (*Category, error) {
	_, executor := getExecutor(dbService, otx...)

	query, args, err := SQLBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "preparing select statement for a category by ID failed")
	}

	var category Category
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
func GetPagedCategories(ctx context.Context, page int, itemsPerPage int, userID int64, dbService *DBService, otx ...*sql.Tx) ([]Category, error) {
	_, executor := getExecutor(dbService, otx...)

	offset := (page - 1) * itemsPerPage

	query, args, err := SQLBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
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

	var categories []Category
	for rows.Next() {
		category := Category{}
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "scanning paginated category row failed")
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// CategoryIDExists checks if a category with the given ID exists in the database.
func CategoryIDExists(ctx context.Context, categoryID int64, userID int64, dbService *DBService, otx ...*sql.Tx) (bool, error) {
	_, executor := getExecutor(dbService, otx...)

	query, args, err := SQLBuilder.Select("1").
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
