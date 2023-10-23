package models

import (
	"errors"
	"log"
)

type Tag struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// PaginationParams holds parameters for paginating database queries
type PaginationParams struct {
	Limit  int
	Offset int
}

// GetAllTags retrieves all tags for a user with pagination.
func GetAllTags(userID string, pagination PaginationParams) ([]Tag, error) {
	rows, err := GetDB().Query("SELECT id, user_id, name FROM tags WHERE user_id=? LIMIT ? OFFSET ?", userID, pagination.Limit, pagination.Offset)
	if err != nil {
		log.Printf("Error querying tags: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name); err != nil {
			log.Printf("Error scanning tag row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// GetTagByID retrieves a tag by its ID.
func GetTagByID(tagID int) (*Tag, error) {
	row := GetDB().QueryRow("SELECT id, user_id, name FROM tags WHERE id=?", tagID)
	var tag Tag
	err := row.Scan(&tag.ID, &tag.UserID, &tag.Name)
	if err != nil {
		log.Printf("Error retrieving tag by ID: %v", err)
		return nil, err
	}
	return &tag, nil
}

// InsertTag adds a new tag to the database.
func InsertTag(tag *Tag) error {
	if len(tag.Name) == 0 {
		return errors.New("tag name cannot be empty")
	}
	if len(tag.Name) > 255 {
		return errors.New("tag name is too long")
	}

	stmt, err := GetDB().Prepare("INSERT INTO tags (user_id, name) VALUES (?, ?)")
	if err != nil {
		log.Printf("Error preparing the statement: %v", err)
		return errors.New("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(tag.UserID, tag.Name)
	if err != nil {
		log.Printf("Error executing the statement with tag %v: %v", tag, err)
		return errors.New("error inserting tag, it may already exist")
	}
	return nil
}

// UpdateTag updates an existing tag in the database.
func UpdateTag(tag *Tag) error {
	if len(tag.Name) == 0 {
		return errors.New("tag name cannot be empty")
	}
	if len(tag.Name) > 255 {
		return errors.New("tag name is too long")
	}

	stmt, err := GetDB().Prepare("UPDATE tags SET user_id=?, name=? WHERE id=?")
	if err != nil {
		log.Printf("Error preparing the statement: %v", err)
		return errors.New("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(tag.UserID, tag.Name, tag.ID)
	if err != nil {
		log.Printf("Error executing the statement with tag %v: %v", tag, err)
		return errors.New("error updating tag")
	}
	return nil
}

// DeleteTag removes a tag from the database.
func DeleteTag(tagID int) error {
	stmt, err := GetDB().Prepare("DELETE FROM tags WHERE id=?")
	if err != nil {
		log.Printf("Error preparing the statement: %v", err)
		return errors.New("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(tagID)
	if err != nil {
		log.Printf("Error executing the statement with tagID %d: %v", tagID, err)
		return errors.New("error deleting tag")
	}
	return nil
}

// GetTagByName retrieves a tag by its name for a specific user.
func GetTagByName(name string, userID string) (*Tag, error) {
	row := GetDB().QueryRow("SELECT id, user_id, name FROM tags WHERE name=? AND user_id=?", name, userID)
	var tag Tag
	err := row.Scan(&tag.ID, &tag.UserID, &tag.Name)
	if err != nil {
		log.Printf("Error retrieving tag by name %s: %v", name, err)
		return nil, err
	}
	return &tag, nil
}
