package models

import (
	"log"
	"time"
	"xspends/util"
)

type TransactionTag struct {
	TransactionID int64     `json:"transaction_id"`
	TagID         int64     `json:"tag_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetTagsByTransactionID retrieves all tags for a specific transaction.
func GetTagsByTransactionID(transactionID int64) ([]Tag, error) {
	rows, err := GetDB().Query("SELECT t.id, t.name FROM tags t JOIN transaction_tags tt ON t.id = tt.tag_id WHERE tt.transaction_id = ?", transactionID)
	if err != nil {
		log.Printf("[ERROR] Querying tags for transaction %d: %v", transactionID, err)
		return nil, util.ErrDatabase
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			log.Printf("[ERROR] Scanning tag row: %v", err)
			return nil, util.ErrDatabase
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// InsertTransactionTag adds a new tag to a specific transaction.
func InsertTransactionTag(transactionID, tagID int64) error {
	_, err := GetDB().Exec("INSERT INTO transaction_tags (transaction_id, tag_id, created_at, updated_at) VALUES (?, ?, ?, ?)", transactionID, tagID, time.Now(), time.Now())
	if err != nil {
		log.Printf("[ERROR] Inserting tag %d for transaction %d: %v", tagID, transactionID, err)
		return util.ErrDatabase
	}
	return nil
}

// DeleteTransactionTag removes a specific tag from a specific transaction.
func DeleteTransactionTag(transactionID, tagID int64) error {
	_, err := GetDB().Exec("DELETE FROM transaction_tags WHERE transaction_id = ? AND tag_id = ?", transactionID, tagID)
	if err != nil {
		log.Printf("[ERROR] Deleting tag %d from transaction %d: %v", tagID, transactionID, err)
		return util.ErrDatabase
	}
	return nil
}

// AddTagsToTransaction adds multiple tags to a specific transaction.
func AddTagsToTransaction(transactionID int64, tags []string, userID int64) error {
	for _, tagName := range tags {
		tag, err := GetTagByName(tagName, userID)
		if err != nil {
			continue
		}
		err = InsertTransactionTag(transactionID, tag.ID)
		if err != nil {
			log.Printf("[ERROR] Associating tag %s with transaction %d: %v", tagName, transactionID, err)
		}
	}
	return nil
}

// UpdateTagsForTransaction updates the tag associations for a specific transaction.
func UpdateTagsForTransaction(transactionID int64, tags []string, userID int64) error {
	if err := RemoveTagsFromTransaction(transactionID); err != nil {
		return err
	}
	return AddTagsToTransaction(transactionID, tags, userID)
}

// RemoveTagsFromTransaction removes all tags associated with a specific transaction.
func RemoveTagsFromTransaction(transactionID int64) error {
	_, err := GetDB().Exec("DELETE FROM transaction_tags WHERE transaction_id = ?", transactionID)
	if err != nil {
		log.Printf("[ERROR] Removing tags from transaction %d: %v", transactionID, err)
		return util.ErrDatabase
	}
	return nil
}
