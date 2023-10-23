package models

import (
	"log"
)

// TransactionTag represents a relationship between a transaction and a tag.
type TransactionTag struct {
	TransactionID int `json:"transaction_id"`
	TagID         int `json:"tag_id"`
}

// GetTagsByTransactionID retrieves all tags for a specific transaction.
func GetTagsByTransactionID(transactionID int) ([]Tag, error) {
	rows, err := GetDB().Query("SELECT tags.id, tags.name FROM tags JOIN transaction_tags ON tags.id = transaction_tags.tag_id WHERE transaction_tags.transaction_id = ?", transactionID)
	if err != nil {
		log.Printf("Error querying tags for transaction %d: %v", transactionID, err)
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			log.Printf("Error scanning tag row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// InsertTransactionTag adds a new tag to a specific transaction.
func InsertTransactionTag(transactionID, tagID int) error {
	_, err := GetDB().Exec("INSERT INTO transaction_tags (transaction_id, tag_id) VALUES (?, ?)", transactionID, tagID)
	if err != nil {
		log.Printf("Error inserting tag %d for transaction %d: %v", tagID, transactionID, err)
		return err
	}
	return nil
}

// DeleteTransactionTag removes a specific tag from a specific transaction.
func DeleteTransactionTag(transactionID, tagID int) error {
	_, err := GetDB().Exec("DELETE FROM transaction_tags WHERE transaction_id = ? AND tag_id = ?", transactionID, tagID)
	if err != nil {
		log.Printf("Error deleting tag %d from transaction %d: %v", tagID, transactionID, err)
		return err
	}
	return nil
}

// AddTagsToTransaction adds multiple tags to a specific transaction.
func AddTagsToTransaction(transactionID int, tags []string, userID string) error {
	for _, tagName := range tags {
		tag, err := GetTagByName(tagName, userID)
		if err != nil {
			continue // if the tag doesn't exist or doesn't belong to the user, skip it
		}
		err = InsertTransactionTag(transactionID, tag.ID)
		if err != nil {
			log.Printf("Error associating tag %s with transaction %d: %v", tagName, transactionID, err)
		}
	}
	return nil
}

// UpdateTagsForTransaction updates the tag associations for a specific transaction.
func UpdateTagsForTransaction(transactionID int, tags []string, userID string) error {
	// First, remove all existing tags for the transaction
	RemoveTagsFromTransaction(transactionID)

	// Then, add the new tags
	return AddTagsToTransaction(transactionID, tags, userID)
}

// RemoveTagsFromTransaction removes all tags associated with a specific transaction.
func RemoveTagsFromTransaction(transactionID int) error {
	_, err := GetDB().Exec("DELETE FROM transaction_tags WHERE transaction_id = ?", transactionID)
	if err != nil {
		log.Printf("Error removing tags from transaction %d: %v", transactionID, err)
		return err
	}
	return nil
}
