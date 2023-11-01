package models

import (
	"database/sql"
	"time"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

const (
	maxCategoryNameLength        = 255
	maxCategoryDescriptionLength = 500
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
func InsertCategory(category *Category) error {
	if category.UserID <= 0 || len(category.Name) == 0 || len(category.Name) > maxCategoryNameLength || len(category.Description) > maxCategoryDescriptionLength {
		return util.ErrInvalidInput
	}

	var err error
	category.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		logrs.WithError(err).Error("Generating Snowflake ID failed")
		return util.ErrDatabase
	}
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	query, args, err := SQLBuilder.Insert("categories").
		Columns("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		Values(category.ID, category.UserID, category.Name, category.Description, category.Icon, category.CreatedAt, category.UpdatedAt).
		ToSql()
	if err != nil {
		logrs.WithError(err).Error("Preparing insert statement failed")
		return util.ErrDatabase
	}

	if _, err := GetDB().Exec(query, args...); err != nil {
		logrs.WithError(err).WithField("category", category).Error("Executing insert statement failed")
		return util.ErrDatabase
	}

	return nil
}

// UpdateCategory updates an existing category in the database.
func UpdateCategory(category *Category) error {
	if category.UserID <= 0 || len(category.Name) == 0 || len(category.Name) > maxCategoryNameLength || len(category.Description) > maxCategoryDescriptionLength {
		return util.ErrInvalidInput
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
		logrs.WithError(err).Error("Preparing update statement failed")
		return util.ErrDatabase
	}

	if _, err := GetDB().Exec(query, args...); err != nil {
		logrs.WithError(err).WithField("category", category).Error("Executing update statement failed")
		return util.ErrDatabase
	}

	return nil
}

// DeleteCategory deletes a category from the database.
func DeleteCategory(categoryID int64, userID int64) error {
	query, args, err := SQLBuilder.Delete("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
		ToSql()
	if err != nil {
		logrs.WithError(err).Error("Preparing delete statement failed")
		return util.ErrDatabase
	}

	if _, err := GetDB().Exec(query, args...); err != nil {
		logrs.WithError(err).WithField("categoryID", categoryID).Error("Executing delete statement failed")
		return util.ErrDatabase
	}

	return nil
}

// GetAllCategories retrieves all categories for a user from the database.
func GetAllCategories(userID int64) ([]Category, error) {
	query, args, err := SQLBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		logrs.WithError(err).Error("Preparing select statement for all categories failed")
		return nil, util.ErrDatabase
	}

	rows, err := GetDB().Query(query, args...)
	if err != nil {
		logrs.WithError(err).Error("Querying categories failed")
		return nil, util.ErrDatabase
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			logrs.WithError(err).Error("Scanning category row failed")
			return nil, util.ErrDatabase
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by its ID for a user from the database.
func GetCategoryByID(categoryID int64, userID int64) (*Category, error) {
	query, args, err := SQLBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
		ToSql()
	if err != nil {
		logrs.WithError(err).Error("Preparing select statement for a category by ID failed")
		return nil, util.ErrDatabase
	}

	var category Category
	err = GetDB().QueryRow(query, args...).Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrs.WithError(err).WithField("categoryID", categoryID).Error("Querying category by ID failed")
		return nil, util.ErrDatabase
	}

	return &category, nil
}

// GetPagedCategories retrieves a paginated list of categories for a user from the database.
func GetPagedCategories(page int, itemsPerPage int, userID int64) ([]Category, error) {
	offset := (page - 1) * itemsPerPage

	query, args, err := SQLBuilder.Select("id", "user_id", "name", "description", "icon", "created_at", "updated_at").
		From("categories").
		Where(squirrel.Eq{"user_id": userID}).
		Limit(uint64(itemsPerPage)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		logrs.WithError(err).Error("Preparing paginated select statement for categories failed")
		return nil, util.ErrDatabase
	}

	rows, err := GetDB().Query(query, args...)
	if err != nil {
		logrs.WithError(err).Error("Querying paginated categories failed")
		return nil, util.ErrDatabase
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			logrs.WithError(err).Error("Scanning paginated category row failed")
			return nil, util.ErrDatabase
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// CategoryIDExists checks if a category with the given ID exists in the database.
func CategoryIDExists(categoryID int64, userID int64) (bool, error) {
	query, args, err := SQLBuilder.Select("1").
		From("categories").
		Where(squirrel.Eq{"id": categoryID, "user_id": userID}).
		Limit(1).
		ToSql()
	if err != nil {
		logrs.WithError(err).Error("Preparing statement to check category existence failed")
		return false, util.ErrDatabase
	}

	var exists int
	err = GetDB().QueryRow(query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logrs.WithError(err).WithFields(logrus.Fields{
			"categoryID": categoryID,
			"userID":     userID,
		}).Error("Checking category existence failed")
		return false, util.ErrDatabase
	}
	return exists == 1, nil
}
