package models

import (
	"errors"
	"log"
)

type Category struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"` // Added this field
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
