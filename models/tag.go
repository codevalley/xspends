package models

import (
	"database/sql"
	"log"
	"time"
	"xspends/util" // Adjust this import to your project's structure
)

const (
	maxTagNameLength = 255
)

type Tag struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginationParams holds parameters for paginating database queries
type PaginationParams struct {
	Limit  int
	Offset int
}

// InsertTag adds a new tag to the database.
func InsertTag(tag *Tag, tx ...*sql.Tx) error {
	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > maxTagNameLength {
		return util.ErrInvalidInput
	}

	var err error
	tag.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		log.Printf("[ERROR] Generating Snowflake ID: %v", err)
		return util.ErrDatabase
	}
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	isExternalTx, txInstance, err := GetTransaction(tx...)
	if err != nil {
		return err
	}

	stmt, err := txInstance.Prepare("INSERT INTO tags (id, user_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[ERROR] Preparing statement: %v", err)
		if !isExternalTx {
			txInstance.Rollback()
		}
		return util.ErrDatabase
	}
	defer stmt.Close()

	_, err = stmt.Exec(tag.ID, tag.UserID, tag.Name, tag.CreatedAt, tag.UpdatedAt)
	if err != nil {
		log.Printf("[ERROR] Executing statement with tag %v: %v", tag, err)
		if !isExternalTx {
			txInstance.Rollback()
		}
		return util.ErrDatabase
	}

	if !isExternalTx {
		return txInstance.Commit()
	}

	return nil
}

// UpdateTag updates an existing tag in the database.
func UpdateTag(tag *Tag, tx ...*sql.Tx) error {
	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > maxTagNameLength {
		return util.ErrInvalidInput
	}

	tag.UpdatedAt = time.Now()

	isExternalTx, txInstance, err := GetTransaction(tx...)
	if err != nil {
		return err
	}

	stmt, err := txInstance.Prepare("UPDATE tags SET name=?, updated_at=? WHERE id=? AND user_id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing statement: %v", err)
		if !isExternalTx {
			txInstance.Rollback()
		}
		return util.ErrDatabase
	}
	defer stmt.Close()

	_, err = stmt.Exec(tag.Name, tag.UpdatedAt, tag.ID, tag.UserID)
	if err != nil {
		log.Printf("[ERROR] Executing statement with tag %v: %v", tag, err)
		if !isExternalTx {
			txInstance.Rollback()
		}
		return util.ErrDatabase
	}

	if !isExternalTx {
		return txInstance.Commit()
	}

	return nil
}

// DeleteTag removes a tag from the database.
func DeleteTag(tagID int64, userID int64) error {
	stmt, err := GetDB().Prepare("DELETE FROM tags WHERE id=? AND user_id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing statement: %v", err)
		return util.ErrDatabase
	}
	defer stmt.Close()

	_, err = stmt.Exec(tagID, userID)
	if err != nil {
		log.Printf("[ERROR] Executing statement with tagID %d: %v", tagID, err)
		return util.ErrDatabase
	}

	return nil
}

// GetTagByID retrieves a tag by its ID.
func GetTagByID(tagID int64, userID int64) (*Tag, error) {
	row := GetDB().QueryRow("SELECT id, user_id, name, created_at, updated_at FROM tags WHERE id=? AND user_id=?", tagID, userID)
	var tag Tag
	err := row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		log.Printf("[ERROR] Retrieving tag by ID: %v", err)
		return nil, err
	}
	return &tag, nil
}

// GetAllTags retrieves all tags for a user with pagination.
func GetAllTags(userID int64, pagination PaginationParams) ([]Tag, error) {
	rows, err := GetDB().Query("SELECT id, user_id, name, created_at, updated_at FROM tags WHERE user_id=? LIMIT ? OFFSET ?", userID, pagination.Limit, pagination.Offset)
	if err != nil {
		log.Printf("[ERROR] Querying tags: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			log.Printf("[ERROR] Scanning tag row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// GetTagByName retrieves a tag by its name for a specific user.
func GetTagByName(name string, userID int64, tx ...*sql.Tx) (*Tag, error) {
	isExternalTx, txInstance, err := GetTransaction(tx...)
	if err != nil {
		return nil, err
	}

	row := txInstance.QueryRow("SELECT id, user_id, name, created_at, updated_at FROM tags WHERE name=? AND user_id=?", name, userID)
	var tag Tag
	err = row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		log.Printf("[ERROR] Retrieving tag by name %s: %v", name, err)
		if !isExternalTx {
			txInstance.Rollback()
		}
		return nil, err
	}

	return &tag, nil
}
