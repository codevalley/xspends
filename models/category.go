package models

import (
	"database/sql"
	"errors"
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

var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrInvalidInput       = errors.New("invalid input data")
	ErrDatabase           = errors.New("database error")
	ErrCategoryNameLength = errors.New("category name length exceeds limit")
	ErrCategoryDescLength = errors.New("category description length exceeds limit")
)

func InsertCategory(category *Category) error {
	if category.UserID <= 0 || len(category.Name) == 0 || len(category.Name) > maxCategoryNameLength || len(category.Description) > maxCategoryDescriptionLength {
		return ErrInvalidInput
	}

	var err error
	category.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		log.Printf("[ERROR] Generating Snowflake ID: %v", err)
		return ErrDatabase // or a more specific error like ErrGeneratingID
	}
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	stmt, err := GetDB().Prepare("INSERT INTO categories (id, user_id, name, description, icon, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[ERROR] Preparing statement: %v", err)
		return ErrDatabase
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.ID, category.UserID, category.Name, category.Description, category.Icon, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		log.Printf("[ERROR] Executing statement with category %v: %v", category, err)
		return ErrDatabase
	}

	return nil
}

func UpdateCategory(category *Category) error {
	if category.UserID <= 0 || len(category.Name) == 0 || len(category.Name) > maxCategoryNameLength || len(category.Description) > maxCategoryDescriptionLength {
		return ErrInvalidInput
	}

	category.UpdatedAt = time.Now()

	stmt, err := GetDB().Prepare("UPDATE categories SET user_id=?, name=?, description=?, icon=?, updated_at=? WHERE id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing statement: %v", err)
		return ErrDatabase
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.UserID, category.Name, category.Description, category.Icon, category.UpdatedAt, category.ID)
	if err != nil {
		log.Printf("[ERROR] Executing statement with category %v: %v", category, err)
		return ErrDatabase
	}

	return nil
}

func DeleteCategory(categoryID int64) error {
	stmt, err := GetDB().Prepare("DELETE FROM categories WHERE id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing statement: %v", err)
		return ErrDatabase
	}
	defer stmt.Close()

	_, err = stmt.Exec(categoryID)
	if err != nil {
		log.Printf("[ERROR] Executing statement with categoryID %d: %v", categoryID, err)
		return ErrDatabase
	}

	return nil
}

func GetAllCategories() ([]Category, error) {
	rows, err := GetDB().Query("SELECT id, user_id, name, description, icon, created_at, updated_at FROM categories")
	if err != nil {
		log.Printf("[ERROR] Querying categories: %v", err)
		return nil, ErrDatabase
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			log.Printf("[ERROR] Scanning category row: %v", err)
			return nil, ErrDatabase
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetCategoryByID(categoryID int64) (*Category, error) {
	row := GetDB().QueryRow("SELECT id, user_id, name, description, icon, created_at, updated_at FROM categories WHERE id=?", categoryID)
	var category Category
	err := row.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrCategoryNotFound
	} else if err != nil {
		log.Printf("[ERROR] Retrieving category by ID %d: %v", categoryID, err)
		return nil, ErrDatabase
	}

	return &category, nil
}

func GetPagedCategories(page, itemsPerPage int) ([]Category, error) {
	offset := (page - 1) * itemsPerPage

	rows, err := GetDB().Query("SELECT id, user_id, name, description, icon, created_at, updated_at FROM categories LIMIT ? OFFSET ?", itemsPerPage, offset)
	if err != nil {
		log.Printf("[ERROR] Querying paginated categories: %v", err)
		return nil, ErrDatabase
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon, &category.CreatedAt, &category.UpdatedAt); err != nil {
			log.Printf("[ERROR] Scanning category row: %v", err)
			return nil, ErrDatabase
		}
		categories = append(categories, category)
	}

	return categories, nil
}
