package models

import (
	"errors"
	"log"
)

type Category struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// InsertCategory adds a new category to the database.
func InsertCategory(category *Category) error {
	// Input validation
	if category.UserID == "" {
		return errors.New("user ID cannot be empty")
	}
	if len(category.Name) == 0 {
		return errors.New("category name cannot be empty")
	}
	if len(category.Name) > 255 {
		return errors.New("category name is too long")
	}
	if len(category.Description) > 500 {
		return errors.New("category description is too long")
	}

	// Prepared statement
	stmt, err := GetDB().Prepare("INSERT INTO categories (id, user_id, name, description, icon) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing the statement: %v", err)
		return errors.New("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.ID, category.UserID, category.Name, category.Description, category.Icon)
	if err != nil {
		log.Printf("Error executing the statement with category %v: %v", category, err)
		return errors.New("error inserting category, it may already exist")
	}

	return nil
}

// UpdateCategory updates an existing category in the database.
func UpdateCategory(category *Category) error {
	// Input validation
	if category.UserID == "" {
		return errors.New("user ID cannot be empty")
	}
	if len(category.Name) == 0 {
		return errors.New("category name cannot be empty")
	}
	if len(category.Name) > 255 {
		return errors.New("category name is too long")
	}
	if len(category.Description) > 500 {
		return errors.New("category description is too long")
	}

	// Prepared statement
	stmt, err := GetDB().Prepare("UPDATE categories SET user_id=?, name=?, description=?, icon=? WHERE id=?")
	if err != nil {
		log.Printf("Error preparing the statement: %v", err)
		return errors.New("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.UserID, category.Name, category.Description, category.Icon, category.ID)
	if err != nil {
		log.Printf("Error executing the statement with category %v: %v", category, err)
		return errors.New("error updating category")
	}

	return nil
}

// DeleteCategory removes a category from the database.
func DeleteCategory(categoryID string) error {
	stmt, err := GetDB().Prepare("DELETE FROM categories WHERE id=?")
	if err != nil {
		log.Printf("Error preparing the statement: %v", err)
		return errors.New("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(categoryID)
	if err != nil {
		log.Printf("Error executing the statement with categoryID %s: %v", categoryID, err)
		return errors.New("error deleting category")
	}

	return nil
}

// GetAllCategories retrieves all categories.
func GetAllCategories() ([]Category, error) {
	rows, err := GetDB().Query("SELECT id, user_id, name, description, icon FROM categories")
	if err != nil {
		log.Printf("Error querying categories: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon); err != nil {
			log.Printf("Error scanning category row: %v", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// GetCategoryByID retrieves a category by its ID.
func GetCategoryByID(categoryID string) (*Category, error) {
	row := GetDB().QueryRow("SELECT id, user_id, name, description, icon FROM categories WHERE id=?", categoryID)
	var category Category
	err := row.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon)
	if err != nil {
		log.Printf("Error retrieving category by ID: %v", err)
		return nil, err
	}
	return &category, nil
}

// GetPagedCategories retrieves categories in a paginated manner.
func GetPagedCategories(page, itemsPerPage int) ([]Category, error) {
	offset := (page - 1) * itemsPerPage

	rows, err := GetDB().Query("SELECT id, user_id, name, description, icon FROM categories LIMIT ? OFFSET ?", itemsPerPage, offset)
	if err != nil {
		log.Printf("Error querying paginated categories: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Description, &category.Icon); err != nil {
			log.Printf("Error scanning category row: %v", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
