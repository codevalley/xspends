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

// InsertCategory inserts a new category into the database.
func InsertCategory(ctx context.Context, category *Category) error {
	if err := validateCategoryInput(category); err != nil {
		return err
	}

	tx, err := GetDB().BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "starting transaction failed")
	}

	category.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "generating Snowflake ID failed")
	}
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	query, args, err := SQLBuilder.Insert("categories").
		Columns("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		Values(category.ID, category.UserID, category.Name, category.Description, category.Icon, category.CreatedAt, category.UpdatedAt).
		ToSql()
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "preparing insert statement failed")
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "preparing SQL statement failed")
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "executing insert statement failed")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "committing transaction failed")
	}

	return nil
}

// UpdateCategory updates an existing category in the database.
func UpdateCategory(ctx context.Context, category *Category) error {
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

	stmt, err := GetDB().PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "preparing SQL statement failed")
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		return errors.Wrap(err, "executing update statement failed")
	}

	return nil
}

func validateCategoryInput(category *Category) error {
	if category.UserID <= 0 || category.Name == "" || len(category.Name) > maxCategoryNameLength || len(category.Description) > maxCategoryDescriptionLength {
		return errors.New(ErrInvalidInput)
	}
	return nil
}

// DeleteCategory deletes a category from the database.
func DeleteCategory(ctx context.Context, categoryID int64, userID int64) error {
	query, args, err := SQLBuilder.Delete("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "preparing delete statement failed")
	}

	if _, err := GetDB().ExecContext(ctx, query, args...); err != nil {
		return errors.Wrap(err, "executing delete statement failed")
	}

	return nil
}

// GetAllCategories retrieves all categories for a user from the database.
func GetAllCategories(ctx context.Context, userID int64) ([]Category, error) {
	query, args, err := SQLBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "preparing select statement for all categories failed")
	}

	rows, err := GetDB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying categories failed")
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		category := &Category{}
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "scanning category row failed")
		}
		categories = append(categories, *category)
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by its ID for a user from the database.
func GetCategoryByID(ctx context.Context, categoryID int64, userID int64) (*Category, error) {
	query, args, err := SQLBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "preparing select statement for a category by ID failed")
	}

	var category Category
	err = GetDB().QueryRowContext(ctx, query, args...).Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found")
		}
		return nil, errors.Wrap(err, "querying category by ID failed")
	}

	return &category, nil
}

// GetPagedCategories retrieves a paginated list of categories for a user from the database.
func GetPagedCategories(ctx context.Context, page int, itemsPerPage int, userID int64) ([]Category, error) {
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

	rows, err := GetDB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying paginated categories failed")
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		category := &Category{}
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "scanning paginated category row failed")
		}
		categories = append(categories, *category)
	}

	return categories, nil
}

// CategoryIDExists checks if a category with the given ID exists in the database.
func CategoryIDExists(ctx context.Context, categoryID int64, userID int64) (bool, error) {
	query, args, err := SQLBuilder.Select("1").
		From("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, errors.Wrap(err, "preparing statement to check category existence failed")
	}

	var exists int
	err = GetDB().QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "checking category existence failed")
	}
	return exists == 1, nil
}
