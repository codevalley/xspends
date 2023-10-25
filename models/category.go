package models

import (
	"database/sql"
	"log"
	"time"
	"xspends/util"
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

func InsertCategory(category *Category) error {
	if category.UserID <= 0 || len(category.Name) == 0 || len(category.Name) > maxCategoryNameLength || len(category.Description) > maxCategoryDescriptionLength {
		return util.ErrInvalidInput
	}

	var err error
	category.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		log.Printf("[ERROR] Generating Snowflake ID: %v", err)
		return util.ErrDatabase // or a more specific error like ErrGeneratingID
	}
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	stmt, err := GetDB().Prepare("INSERT INTO categories (id, user_id, name, description, icon, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[ERROR] Preparing statement: %v", err)
		return util.ErrDatabase
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.ID, category.UserID, category.Name, category.Description, category.Icon, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		log.Printf("[ERROR] Executing statement with category %v: %v", category, err)
		return util.ErrDatabase
	}

	return nil
}

func UpdateCategory(category *Category) error {
	if category.UserID <= 0 || len(category.Name) == 0 || len(category.Name) > maxCategoryNameLength || len(category.Description) > maxCategoryDescriptionLength {
		return util.ErrInvalidInput
	}

	category.UpdatedAt = time.Now()

	stmt, err := GetDB().Prepare("UPDATE categories SET name=?, description=?, icon=?, updated_at=? WHERE id=? AND user_id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing statement: %v", err)
		return util.ErrDatabase
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.Name, category.Description, category.Icon, category.UpdatedAt, category.ID, category.UserID)
	if err != nil {
		log.Printf("[ERROR] Executing statement with category %v: %v", category, err)
		return util.ErrDatabase
	}

	return nil
}

func DeleteCategory(categoryID int64, userID int64) error {
	stmt, err := GetDB().Prepare("DELETE FROM categories WHERE id=? AND user_id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing statement: %v", err)
		return util.ErrDatabase
	}
	defer stmt.Close()

	_, err = stmt.Exec(categoryID)
	if err != nil {
		log.Printf("[ERROR] Executing statement with categoryID %d: %v", categoryID, err)
		return util.ErrDatabase
	}

	return nil
}

func GetAllCategories(userID int64) ([]Category, error) {
	rows, err := GetDB().Query("SELECT id, user_id, name, description, icon, created_at, updated_at FROM categories WHERE user_id=?", userID)
	if err != nil {
		log.Printf("[ERROR] Querying categories: %v", err)
		return nil, util.ErrDatabase
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			log.Printf("[ERROR] Scanning category row: %v", err)
			return nil, util.ErrDatabase
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetCategoryByID(categoryID int64, userID int64) (*Category, error) {
	row := GetDB().QueryRow("SELECT id, user_id, name, description, icon, created_at, updated_at FROM categories WHERE id=? AND user_id=?", categoryID, userID)
	var category Category
	err := row.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, util.ErrCategoryNotFound
	} else if err != nil {
		log.Printf("[ERROR] Retrieving category by ID %d: %v", categoryID, err)
		return nil, util.ErrDatabase
	}

	return &category, nil
}

func GetPagedCategories(page int, itemsPerPage int, userID int64) ([]Category, error) {
	offset := (page - 1) * itemsPerPage

	rows, err := GetDB().Query("SELECT id, user_id, name, description, icon, created_at, updated_at FROM categories WHERE user_id=? LIMIT ? OFFSET ?", userID, itemsPerPage, offset)
	if err != nil {
		log.Printf("[ERROR] Querying paginated categories: %v", err)
		return nil, util.ErrDatabase
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			log.Printf("[ERROR] Scanning category row: %v", err)
			return nil, util.ErrDatabase
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// CategoryIDExists checks if a category with the given ID exists in the database.
func CategoryIDExists(categoryID int64, userID int64) (bool, error) {
	stmt, err := GetDB().Prepare("SELECT 1 FROM categories WHERE id=? and user_id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing select statement to check if category exists by ID: %v", err)
		return false, err
	}
	defer stmt.Close()

	var exists int
	err = stmt.QueryRow(categoryID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("[ERROR] Checking if category exists by ID: %v", err)
		return false, err
	}
	return exists == 1, nil
}
