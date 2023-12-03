/*
MIT License

Copyright (c) 2023 Narayan Babu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package models

import (
	"context"
	"database/sql"

	"time"

	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

const (
	TransactionTypeIncome  = "INCOME"
	TransactionTypeExpense = "EXPENSE"
	SortOrderAsc           = "ASC"
	SortOrderDesc          = "DESC"
)

type Transaction struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	SourceID    int64     `json:"source_id"`
	Tags        []string  `json:"tags"`
	CategoryID  int64     `json:"category_id"`
	Timestamp   time.Time `json:"timestamp"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
}

type TransactionFilter struct {
	UserID       int64
	StartDate    string
	EndDate      string
	Tags         []string
	Category     string
	Type         string
	Description  string
	MinAmount    float64
	MaxAmount    float64
	SortBy       string
	SortOrder    string // "ASC" or "DESC"
	Page         int
	ItemsPerPage int
}

// InsertTransaction inserts a new transaction into the database.// InsertTransaction inserts a new transaction into the database.
func InsertTransaction(ctx context.Context, txn Transaction, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	txn.ID, _ = util.GenerateSnowflakeID()
	txn.Timestamp = time.Now()

	if err := validateForeignKeyReferences(ctx, txn, otx...); err != nil {
		return errors.Wrap(err, "validating foreign key references failed")
	}

	query, args, err := squirrel.Insert("transactions").
		Columns("id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description").
		Values(txn.ID, txn.UserID, txn.SourceID, txn.CategoryID, txn.Timestamp, txn.Amount, txn.Type, txn.Description).
		PlaceholderFormat(squirrel.Question).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query for transaction")
	}

	if _, err := executor.ExecContext(ctx, query, args...); err != nil {
		return errors.Wrap(err, "insert transaction failed")
	}

	if err := addMissingTags(ctx, txn.ID, txn.Tags, txn.UserID, otx...); err != nil {
		return errors.Wrap(err, "handling transaction tags failed")
	}
	// Associate tags with the transaction
	if err := AddTagsToTransaction(ctx, txn.ID, txn.Tags, txn.UserID, otx...); err != nil {
		return errors.Wrap(err, "adding tags to transaction failed")
	}
	commitOrRollback(executor, isExternalTx, err)
	return nil
}

// UpdateTransaction updates an existing transaction in the database.
func UpdateTransaction(ctx context.Context, txn Transaction, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	// Validate foreign key references
	if err := validateForeignKeyReferences(ctx, txn, otx...); err != nil {
		return errors.Wrap(err, "validating foreign key references failed")
	}

	// Update transaction in the database
	query, args, err := SQLBuilder.Update("transactions").
		Set("source_id", txn.SourceID).
		Set("category_id", txn.CategoryID).
		Set("amount", txn.Amount).
		Set("type", txn.Type).
		Set("description", txn.Description).
		Where(squirrel.Eq{"id": txn.ID, "user_id": txn.UserID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build update query for transaction")
	}

	if _, err := executor.ExecContext(ctx, query, args...); err != nil {
		return errors.Wrap(err, "update transaction failed")
	}

	// Add any missing tags and update tags associated with the transaction
	if err := addMissingTags(ctx, txn.ID, txn.Tags, txn.UserID, otx...); err != nil {
		return errors.Wrap(err, "adding missing tags failed")
	}
	if err := UpdateTagsForTransaction(ctx, txn.ID, txn.Tags, txn.UserID, otx...); err != nil {
		return errors.Wrap(err, "updating tags for transaction failed")
	}

	commitOrRollback(executor, isExternalTx, err)

	return nil
}

// DeleteTransaction removes a transaction from the database.
func DeleteTransaction(ctx context.Context, transactionID int64, userID int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	query, args, err := SQLBuilder.Delete("transactions").
		Where(squirrel.Eq{"id": transactionID, "user_id": userID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build delete query for transaction")
	}

	if _, err := executor.ExecContext(ctx, query, args...); err != nil {
		return errors.Wrap(err, "delete transaction failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	// if !isExternalTx {
	// 	if tx, ok := executor.(*sql.Tx); ok {
	// 		if err := tx.Commit(); err != nil {
	// 			tx.Rollback()
	// 			return errors.Wrap(err, "committing transaction failed")
	// 		}
	// 	}
	// }

	return nil
}

// GetTransactionByID retrieves a single transaction from the database by its ID.
func GetTransactionByID(ctx context.Context, transactionID int64, userID int64, otx ...*sql.Tx) (*Transaction, error) {
	_, executor := getExecutor(otx...)

	query, args, err := SQLBuilder.Select("id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description").
		From("transactions").
		Where(squirrel.Eq{"id": transactionID, "user_id": userID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query for retrieving transaction by ID")
	}

	row := executor.QueryRowContext(ctx, query, args...)
	var transaction Transaction
	if err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.SourceID, &transaction.CategoryID, &transaction.Timestamp, &transaction.Amount, &transaction.Type, &transaction.Description); err != nil {
		return nil, errors.Wrap(err, "get transaction by ID failed")
	}

	getTagsForTransaction(ctx, &transaction, otx...)

	return &transaction, nil
}

// GetTransactionsByFilter retrieves a list of transactions from the database based on a set of filters.
func GetTransactionsByFilter(ctx context.Context, filter TransactionFilter, otx ...*sql.Tx) ([]Transaction, error) {
	_, executor := getExecutor(otx...)
	query := SQLBuilder.Select("id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description").
		From("transactions").
		Where(squirrel.Eq{"user_id": filter.UserID})

	if filter.StartDate != "" {
		query = query.Where("timestamp >= ?", filter.StartDate)
	}

	if filter.EndDate != "" {
		query = query.Where("timestamp <= ?", filter.EndDate)
	}

	if filter.Category != "" {
		query = query.Where("category_id = ?", filter.Category)
	}

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.Description != "" {
		query = query.Where("description LIKE ?", "%"+filter.Description+"%")
	}

	if len(filter.Tags) > 0 {
		tagsSubQuery := SQLBuilder.Select("transaction_id").
			From("transaction_tags").
			Where("tag_id IN ?", filter.Tags)
		query = query.Where("id IN ?", tagsSubQuery)
	}

	if filter.MinAmount > 0 {
		query = query.Where("amount >= ?", filter.MinAmount)
	}

	if filter.MaxAmount > 0 {
		query = query.Where("amount <= ?", filter.MaxAmount)
	}

	if filter.SortBy != "" {
		order := "ASC"
		if filter.SortOrder == SortOrderDesc {
			order = "DESC"
		}
		query = query.OrderBy(filter.SortBy + " " + order)
	}

	if filter.Page > 0 && filter.ItemsPerPage > 0 {
		offset := uint64((filter.Page - 1) * filter.ItemsPerPage)
		query = query.Offset(offset).Limit(uint64(filter.ItemsPerPage))
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "constructing SQL query failed")
	}

	rows, err := executor.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying transactions by filter failed")
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.SourceID, &transaction.CategoryID, &transaction.Timestamp, &transaction.Amount, &transaction.Type, &transaction.Description); err != nil {
			return nil, errors.Wrap(err, "scanning transaction failed")
		}
		getTagsForTransaction(ctx, &transaction, otx...)
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "processing rows failed")
	}

	return transactions, nil
}

func getTagsForTransaction(ctx context.Context, transaction *Transaction, otx ...*sql.Tx) error {
	tags, err := GetTagsByTransactionID(ctx, transaction.ID, otx...)
	if err != nil {
		return errors.Wrap(err, "Couldn't fetch tags for the transaction")
	}

	transaction.Tags = make([]string, len(tags))
	for i, tag := range tags {
		transaction.Tags[i] = tag.Name
	}

	return nil
}

// validateForeignKeyReferences checks if the foreign keys in the transaction exist.
func validateForeignKeyReferences(ctx context.Context, txn Transaction, otx ...*sql.Tx) error {
	// Check if the user exists
	userExists, err := UserIDExists(ctx, txn.UserID, otx...)
	if err != nil {
		return errors.Wrap(err, "error checking if user exists")
	}
	if !userExists {
		return errors.New("user does not exist")
	}

	// Check if the source exists
	sourceExists, err := SourceIDExists(ctx, txn.SourceID, txn.UserID, otx...)
	if err != nil {
		return errors.Wrap(err, "error checking if source exists")
	}
	if !sourceExists {
		return errors.New("source does not exist")
	}

	// Check if the category exists
	categoryExists, err := CategoryIDExists(ctx, txn.CategoryID, txn.UserID, otx...)
	if err != nil {
		return errors.Wrap(err, "error checking if category exists")
	}
	if !categoryExists {
		return errors.New("category does not exist")
	}

	return nil
}

// addMissingTags ensures that all tags are present in the database and associates them with the user.
func addMissingTags(ctx context.Context, transactionID int64, tagNames []string, userID int64, otx ...*sql.Tx) error {
	// Ensure all tags are present in the database
	for _, tagName := range tagNames {
		tag, _ := GetTagByName(ctx, tagName, userID, otx...)

		if tag == nil {
			// Tag does not exist; create it
			newTag := Tag{
				UserID: userID,
				Name:   tagName,
			}
			if err := InsertTag(ctx, &newTag, otx...); err != nil {
				return errors.Wrapf(err, "failed to insert new tag '%s'", tagName)
			}
		}
	}

	return nil
}
