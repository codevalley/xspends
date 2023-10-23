package models

import (
	"database/sql"
	"log"
	"time"
)

const (
	TransactionTypeIncome  = "INCOME"
	TransactionTypeExpense = "EXPENSE"
)

type Transaction struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	SourceID   string    `json:"source_id"`
	Tags       []string  `json:"tags"`
	CategoryID string    `json:"category_id"`
	Timestamp  time.Time `json:"timestamp"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
}

func InsertTransaction(transaction Transaction) error {
	tx, err := GetDB().Begin()
	if err != nil {
		return err
	}

	// Insert the transaction
	_, err = tx.Exec("INSERT INTO transactions (id, user_id, source_id, category_id, timestamp, amount, type) VALUES (?, ?, ?, ?, ?, ?, ?)", transaction.ID, transaction.UserID, transaction.SourceID, transaction.CategoryID, transaction.Timestamp, transaction.Amount, transaction.Type)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Handle tags
	for _, tagName := range transaction.Tags {
		var tagID int
		// Check if the tag exists in the `tags` table
		err := tx.QueryRow("SELECT id FROM tags WHERE name=?", tagName).Scan(&tagID)
		if err == sql.ErrNoRows {
			// If the tag doesn't exist, insert it
			res, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tagName)
			if err != nil {
				tx.Rollback()
				return err
			}
			// Get the ID of the newly inserted tag
			lastID, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				return err
			}
			tagID = int(lastID)
		} else if err != nil {
			tx.Rollback()
			return err
		}

		// Insert into the `transaction_tags` table
		_, err = tx.Exec("INSERT INTO transaction_tags (transaction_id, tag_id) VALUES (?, ?)", transaction.ID, tagID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func UpdateTransaction(transaction Transaction) error {
	tx, err := GetDB().Begin()
	if err != nil {
		return err
	}

	// Update the transaction
	_, err = tx.Exec("UPDATE transactions SET user_id=?, source_id=?, category_id=?, timestamp=?, amount=?, type=? WHERE id=?", transaction.UserID, transaction.SourceID, transaction.CategoryID, transaction.Timestamp, transaction.Amount, transaction.Type, transaction.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// First, delete all existing tag associations for the transaction
	_, err = tx.Exec("DELETE FROM transaction_tags WHERE transaction_id=?", transaction.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then, re-insert the tags
	for _, tagName := range transaction.Tags {
		var tagID int
		// Check if the tag exists in the `tags` table
		err := tx.QueryRow("SELECT id FROM tags WHERE name=?", tagName).Scan(&tagID)
		if err == sql.ErrNoRows {
			// If the tag doesn't exist, insert it
			res, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tagName)
			if err != nil {
				tx.Rollback()
				return err
			}
			// Get the ID of the newly inserted tag
			lastID, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				return err
			}
			tagID = int(lastID)
		} else if err != nil {
			tx.Rollback()
			return err
		}

		// Insert into the `transaction_tags` table
		_, err = tx.Exec("INSERT INTO transaction_tags (transaction_id, tag_id) VALUES (?, ?)", transaction.ID, tagID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func DeleteTransaction(transactionID string) error {
	_, err := GetDB().Exec("DELETE FROM transactions WHERE id=?", transactionID)
	if err != nil {
		log.Println("Error deleting transaction:", err)
		return err
	}
	return nil
}

func GetTransactionByID(transactionID string) (*Transaction, error) {
	row := GetDB().QueryRow("SELECT id, user_id, source_id, tags, category_id, timestamp, amount, type FROM transactions WHERE id=?", transactionID)
	var transaction Transaction
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.SourceID, &transaction.Tags, &transaction.CategoryID, &transaction.Timestamp, &transaction.Amount, &transaction.Type)
	if err != nil {
		log.Println("Error retrieving transaction by ID:", err)
		return nil, err
	}
	return &transaction, nil
}
