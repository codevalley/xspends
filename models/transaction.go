package models

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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

type TransactionFilter struct {
	UserID       string
	StartDate    string
	EndDate      string
	Tags         []string
	Category     string
	Type         string
	MinAmount    float64
	MaxAmount    float64
	SortBy       string
	SortOrder    string // "ASC" or "DESC"
	Page         int
	ItemsPerPage int
}

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidFilter       = errors.New("invalid filter provided")
)

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
	_, err = tx.Exec("UPDATE transactions SET user_id=?, source_id=?, category_id=?, timestamp=?, amount=?, type=? WHERE id=? AND user_id=?", transaction.UserID, transaction.SourceID, transaction.CategoryID, transaction.Timestamp, transaction.Amount, transaction.Type, transaction.ID, transaction.UserID)
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

func DeleteTransaction(transactionID string, userID string) error {
	_, err := GetDB().Exec("DELETE FROM transactions WHERE id=? AND user_id=?", transactionID, userID)
	if err != nil {
		log.Println("Error deleting transaction:", err)
		return err
	}
	return nil
}

func GetTransactionByID(transactionID string, userID string) (*Transaction, error) {
	row := GetDB().QueryRow("SELECT id, user_id, source_id, category_id, timestamp, amount, type FROM transactions WHERE id=? AND user_id=?", transactionID, userID)
	var transaction Transaction
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.SourceID, &transaction.CategoryID, &transaction.Timestamp, &transaction.Amount, &transaction.Type)
	if err != nil {
		log.Println("Error retrieving transaction by ID:", err)
		return nil, err
	}
	return &transaction, nil
}

func ConstructQuery(filter TransactionFilter) (string, []interface{}, error) {
	var queryBuffer bytes.Buffer
	var args []interface{}
	var conditions []string

	// Always filter by user ID
	conditions = append(conditions, "user_id = ?")
	args = append(args, filter.UserID)

	if filter.StartDate != "" {
		// Optional: Check if the date format is correct here
		conditions = append(conditions, "timestamp >= ?")
		args = append(args, filter.StartDate)
	}

	if filter.EndDate != "" {
		// Optional: Check if the date format is correct here
		conditions = append(conditions, "timestamp <= ?")
		args = append(args, filter.EndDate)
	}

	if filter.Category != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, filter.Category)
	}

	if filter.Type != "" {
		if filter.Type != TransactionTypeIncome && filter.Type != TransactionTypeExpense {
			return "", nil, fmt.Errorf("invalid transaction type provided")
		}
		conditions = append(conditions, "type = ?")
		args = append(args, filter.Type)
	}

	// Handling tags is a bit more involved, so we'll add the logic here
	if len(filter.Tags) > 0 {
		tagsPlaceholder := strings.Repeat("?,", len(filter.Tags)-1) + "?"
		conditions = append(conditions, fmt.Sprintf("id IN (SELECT transaction_id FROM transaction_tags WHERE tag_id IN (%s))", tagsPlaceholder))
		for _, tag := range filter.Tags {
			args = append(args, tag)
		}
	}

	if filter.MinAmount > 0 {
		conditions = append(conditions, "amount >= ?")
		args = append(args, filter.MinAmount)
	}

	if filter.MaxAmount > 0 {
		conditions = append(conditions, "amount <= ?")
		args = append(args, filter.MaxAmount)
	}

	// Sort
	if filter.SortBy != "" {
		// Validate the sortBy field here
		allowedSortFields := []string{"timestamp", "amount"} // Example allowed fields
		if !contains(allowedSortFields, filter.SortBy) {
			return "", nil, fmt.Errorf("invalid sort field provided")
		}

		sortDirection := "ASC"
		if filter.SortOrder == "DESC" {
			sortDirection = "DESC"
		}
		queryBuffer.WriteString(fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, sortDirection))
	}

	// Pagination
	if filter.ItemsPerPage > 0 {
		queryBuffer.WriteString(" LIMIT ?")
		args = append(args, filter.ItemsPerPage)
		if filter.Page > 1 {
			offset := (filter.Page - 1) * filter.ItemsPerPage
			queryBuffer.WriteString(" OFFSET ?")
			args = append(args, offset)
		}
	}

	return queryBuffer.String(), args, nil
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func GetTransactionsByFilter(filter TransactionFilter) ([]Transaction, error) {
	query, args, err := ConstructQuery(filter)
	if err != nil {
		log.Printf("Error constructing query: %v", err)
		return nil, err
	}

	rows, err := GetDB().Query(query, args...)
	if err != nil {
		log.Printf("Error querying transactions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.SourceID, &transaction.CategoryID, &transaction.Timestamp, &transaction.Amount, &transaction.Type); err != nil {
			log.Printf("Error scanning transaction row: %v", err)
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		return nil, ErrTransactionNotFound
	}

	return transactions, rows.Err()
}
