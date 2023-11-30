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

// InsertTransaction inserts a new transaction into the database.
func InsertTransaction(ctx context.Context, txn Transaction, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTxn(ctx, otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	txn.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		return errors.Wrap(err, "generating Snowflake ID failed")
	}

	err = validateForeignKeyReferences(ctx, txn, tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "validating foreign key references failed")
	}

	query := SQLBuilder.Insert("transactions").
		Columns("id", "user_id", "source_id", "category_id", "amount", "type", "description").
		Values(txn.ID, txn.UserID, txn.SourceID, txn.CategoryID, txn.Amount, txn.Type, txn.Description)

	_, err = query.RunWith(tx).ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "insert transaction failed")
	}

	err = addMissingTags(ctx, txn.Tags, txn.UserID, tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "adding missing tags failed")
	}

	err = AddTagsToTransaction(ctx, txn.ID, txn.Tags, txn.UserID, tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "adding tags to transaction failed")
	}

	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}

// UpdateTransaction updates an existing transaction in the database.
func UpdateTransaction(ctx context.Context, txn Transaction, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTxn(ctx, otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	err = validateForeignKeyReferences(ctx, txn, tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "validating foreign key references failed")
	}

	query := SQLBuilder.Update("transactions").
		Set("source_id", txn.SourceID).
		Set("category_id", txn.CategoryID).
		Set("amount", txn.Amount).
		Set("type", txn.Type).
		Set("description", txn.Description).
		Where(squirrel.Eq{"id": txn.ID, "user_id": txn.UserID})

	_, err = query.RunWith(tx).ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update transaction failed")
	}

	err = addMissingTags(ctx, txn.Tags, txn.UserID, tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "adding missing tags failed")
	}

	err = UpdateTagsForTransaction(ctx, txn.ID, txn.Tags, txn.UserID, tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "updating tags for transaction failed")
	}

	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}

// DeleteTransaction removes a transaction from the database.
func DeleteTransaction(ctx context.Context, transactionID int64, userID int64) error {
	query := SQLBuilder.Delete("transactions").Where(squirrel.Eq{"id": transactionID, "user_id": userID})

	_, err := query.RunWith(GetDB()).ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "delete transaction failed")
	}
	return nil
}

// GetTransactionByID retrieves a single transaction from the database by its ID.
func GetTransactionByID(ctx context.Context, transactionID int64, userID int64) (*Transaction, error) {
	query := SQLBuilder.Select("id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description").
		From("transactions").
		Where(squirrel.Eq{"id": transactionID, "user_id": userID})

	row := query.RunWith(GetDB()).QueryRowContext(ctx)
	var transaction Transaction
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.SourceID, &transaction.CategoryID, &transaction.Timestamp, &transaction.Amount, &transaction.Type, &transaction.Description)
	if err != nil {
		return nil, errors.Wrap(err, "get transaction by ID failed")
	}
	return &transaction, nil
}

// GetTransactionsByFilter retrieves a list of transactions from the database based on a set of filters.
func GetTransactionsByFilter(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
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

	rows, err := GetDB().QueryContext(ctx, sql, args...)
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
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "processing rows failed")
	}

	return transactions, nil
}

// validateForeignKeyReferences checks if the foreign keys in the transaction exist.
func validateForeignKeyReferences(ctx context.Context, transaction Transaction, tx *sql.Tx) error {
	userExists, userErr := UserIDExists(ctx, transaction.UserID, tx)
	sourceExists, sourceErr := SourceIDExists(ctx, transaction.SourceID, transaction.UserID, GetDBService())
	categoryExists, categoryErr := CategoryIDExists(ctx, transaction.CategoryID, transaction.UserID)

	// Return an error if there was a problem checking any reference
	if userErr != nil {
		return errors.Wrap(userErr, "error checking if user exists")
	}
	if sourceErr != nil {
		return errors.Wrap(sourceErr, "error checking if source exists")
	}
	if categoryErr != nil {
		return errors.Wrap(categoryErr, "error checking if category exists")
	}

	// If any of the references do not exist, return an error
	if !userExists || !sourceExists || !categoryExists {
		return errors.New("invalid foreign key references")
	}

	return nil
}

// addMissingTags ensures that all tags are present in the database and associates them with the user.
func addMissingTags(ctx context.Context, tags []string, userID int64, tx ...*sql.Tx) error {
	isExternalTx, txInstance, err := GetTxn(ctx, tx...)
	if err != nil {
		return errors.Wrap(err, "error obtaining transaction for adding tags")
	}

	// Handle tags
	for _, tagName := range tags {
		tag, _ := GetTagByName(ctx, tagName, userID, txInstance)
		if tag == nil {
			// The tag does not exist, so create it
			var nTag Tag
			nTag.Name = tagName
			nTag.UserID = userID
			err = InsertTag(ctx, &nTag, txInstance)
			if err != nil {
				if !isExternalTx {
					txInstance.Rollback()
				}
				return errors.Wrap(err, "error inserting new tag")
			}
		}
	}

	// Only commit if this function created the transaction
	if !isExternalTx {
		err = txInstance.Commit()
		if err != nil {
			return errors.Wrap(err, "error committing transaction for adding tags")
		}
	}

	return nil
}
