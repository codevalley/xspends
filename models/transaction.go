package models

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"xspends/util"
)

const (
	TransactionTypeIncome  = "INCOME"
	TransactionTypeExpense = "EXPENSE"
	SortOrderAsc           = "ASC"
	SortOrderDesc          = "DESC"
)

type Transaction struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	SourceID   int64     `json:"source_id"`
	Tags       []string  `json:"tags"`
	CategoryID int64     `json:"category_id"`
	Timestamp  time.Time `json:"timestamp"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
}

type TransactionFilter struct {
	UserID       int64
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
	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	var err1 error
	transaction.ID, err = util.GenerateSnowflakeID()
	if err1 != nil {
		log.Printf("[ERROR] Generating Snowflake ID: %v", err)
		return ErrDatabase // or a more specific error like ErrGeneratingID
	}

	err = validateForeignKeyReferences(transaction)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO transactions (id, user_id, source_id, category_id, amount, type) VALUES (?,?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(transaction.ID, transaction.UserID, transaction.SourceID, transaction.CategoryID, transaction.Amount, transaction.Type)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = AddTagsToTransaction(transaction.ID, transaction.Tags, transaction.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func UpdateTransaction(transaction Transaction) error {
	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = validateForeignKeyReferences(transaction)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare("UPDATE transactions SET user_id=?, source_id=?, category_id=?, timestamp=?, amount=?, type=? WHERE id=? AND user_id=?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(transaction.UserID, transaction.SourceID, transaction.CategoryID, transaction.Timestamp, transaction.Amount, transaction.Type, transaction.ID, transaction.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = UpdateTagsForTransaction(transaction.ID, transaction.Tags, transaction.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DeleteTransaction(transactionID int64, userID int64) error {
	_, err := GetDB().Exec("DELETE FROM transactions WHERE id=? AND user_id=?", transactionID, userID)
	if err != nil {
		log.Println("Error deleting transaction:", err)
		return err
	}
	return nil
}

func GetTransactionByID(transactionID int64, userID int64) (*Transaction, error) {
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
		conditions = append(conditions, "timestamp >= ?")
		args = append(args, filter.StartDate)
	}

	if filter.EndDate != "" {
		conditions = append(conditions, "timestamp <= ?")
		args = append(args, filter.EndDate)
	}

	if filter.Category != "" {
		conditions = append(conditions, "category_id = ?")
		args = append(args, filter.Category)
	}

	if filter.Type != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, filter.Type)
	}

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

	// Combine all conditions with 'AND'
	combinedConditions := strings.Join(conditions, " AND ")

	// Base SQL query
	queryBuffer.WriteString("SELECT * FROM transactions ")

	if len(conditions) > 0 {
		queryBuffer.WriteString(" WHERE ")
		queryBuffer.WriteString(combinedConditions)
	}

	// Sort
	if filter.SortBy != "" {
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

func validateForeignKeyReferences(transaction Transaction) error {
	userExists, userErr := UserIDExists(transaction.UserID)
	sourceExists, sourceErr := SourceIDExists(transaction.SourceID, transaction.UserID)
	categoryExists, categoryErr := CategoryIDExists(transaction.CategoryID, transaction.UserID)

	if userErr != nil || sourceErr != nil || categoryErr != nil {
		return errors.New("error checking foreign key references")
	}

	if !userExists || !sourceExists || !categoryExists {
		return errors.New("invalid foreign key references")
	}

	return nil
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
		log.Printf("Error querying transactions: %v \n %s", err, query)
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
