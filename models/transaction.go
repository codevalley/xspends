package models

import (
	"database/sql"
	"errors"
	"time"

	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
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
func InsertTransaction(txn Transaction, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return err
	}

	txn.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		logrs.WithError(err).Error("Generating Snowflake ID failed")
		return util.ErrDatabase // or a more specific error like ErrGeneratingID
	}

	err = validateForeignKeyReferences(txn)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := SQLBuilder.Insert("transactions").
		Columns("id", "user_id", "source_id", "category_id", "amount", "type", "description").
		Values(txn.ID, txn.UserID, txn.SourceID, txn.CategoryID, txn.Amount, txn.Type, txn.Description)

	_, err = query.RunWith(tx).Exec()
	if err != nil {
		tx.Rollback()
		logrs.WithError(err).Error("Insert transaction failed")
		return err
	}

	err = addMissingTags(txn.Tags, txn.UserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = AddTagsToTransaction(txn.ID, txn.Tags, txn.UserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !isExternalTx {
		return tx.Commit()
	}
	return nil
}

// UpdateTransaction updates an existing transaction in the database.
func UpdateTransaction(txn Transaction, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return err
	}

	err = validateForeignKeyReferences(txn)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := SQLBuilder.Update("transactions").
		Set("source_id", txn.SourceID).
		Set("category_id", txn.CategoryID).
		Set("amount", txn.Amount).
		Set("type", txn.Type).
		Set("description", txn.Description).
		Where(squirrel.Eq{"id": txn.ID, "user_id": txn.UserID})

	_, err = query.RunWith(tx).Exec()
	if err != nil {
		tx.Rollback()
		logrs.WithError(err).WithField("transaction_id", txn.ID).Error("Update transaction failed")
		return err
	}

	err = addMissingTags(txn.Tags, txn.UserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = UpdateTagsForTransaction(txn.ID, txn.Tags, txn.UserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !isExternalTx {
		return tx.Commit()
	}
	return nil
}

// DeleteTransaction removes a transaction from the database.
func DeleteTransaction(transactionID int64, userID int64) error {
	query := SQLBuilder.Delete("transactions").Where(squirrel.Eq{"id": transactionID, "user_id": userID})

	_, err := query.RunWith(GetDB()).Exec()
	if err != nil {
		logrs.WithError(err).WithFields(logrus.Fields{
			"transaction_id": transactionID,
			"user_id":        userID,
		}).Error("Delete transaction failed")
		return err
	}
	return nil
}

// GetTransactionByID retrieves a single transaction from the database by its ID.
func GetTransactionByID(transactionID int64, userID int64) (*Transaction, error) {
	query := SQLBuilder.Select("id", "user_id", "source_id", "category_id", "timestamp", "amount", "type", "description").
		From("transactions").
		Where(squirrel.Eq{"id": transactionID, "user_id": userID})

	row := query.RunWith(GetDB()).QueryRow()
	var transaction Transaction
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.SourceID, &transaction.CategoryID, &transaction.Timestamp, &transaction.Amount, &transaction.Type, &transaction.Description)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionsByFilter retrieves a list of transactions from the database based on a set of filters.
func GetTransactionsByFilter(filter TransactionFilter) ([]Transaction, error) {
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
		logrs.WithError(err).Error("Constructing SQL query failed")
		return nil, err
	}

	rows, err := GetDB().Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.SourceID, &transaction.CategoryID, &transaction.Timestamp, &transaction.Amount, &transaction.Type, &transaction.Description); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// validateForeignKeyReferences checks if the foreign keys in the transaction exist.
func validateForeignKeyReferences(transaction Transaction) error {
	userExists, userErr := UserIDExists(transaction.UserID)
	sourceExists, sourceErr := SourceIDExists(transaction.SourceID, transaction.UserID)
	categoryExists, categoryErr := CategoryIDExists(transaction.CategoryID, transaction.UserID)

	// Log and return an error if there was a problem checking any reference
	if userErr != nil {
		logrs.WithError(userErr).WithField("user_id", transaction.UserID).Error("Error checking if user exists")
		return userErr
	}
	if sourceErr != nil {
		logrs.WithError(sourceErr).WithField("source_id", transaction.SourceID).Error("Error checking if source exists")
		return sourceErr
	}
	if categoryErr != nil {
		logrs.WithError(categoryErr).WithField("category_id", transaction.CategoryID).Error("Error checking if category exists")
		return categoryErr
	}

	// If any of the references do not exist, log and return an error
	if !userExists || !sourceExists || !categoryExists {
		logrs.WithFields(logrus.Fields{
			"user_exists":     userExists,
			"source_exists":   sourceExists,
			"category_exists": categoryExists,
		}).Error("Invalid foreign key references")
		return errors.New("invalid foreign key references")
	}

	return nil
}

// addMissingTags ensures that all tags are present in the database and associates them with the user.
func addMissingTags(tags []string, userID int64, tx ...*sql.Tx) error {
	isExternalTx, txInstance, err := GetTransaction(tx...)
	if err != nil {
		logrs.WithError(err).Error("Error obtaining transaction for adding tags")
		return err
	}

	// Handle tags
	for _, tagName := range tags {
		tag, _ := GetTagByName(tagName, userID, txInstance)
		if tag == nil {
			// The tag does not exist, so create it
			var nTag Tag
			nTag.Name = tagName
			nTag.UserID = userID
			err = InsertTag(&nTag, txInstance)
			if err != nil {
				logrs.WithError(err).WithFields(logrus.Fields{
					"tag":     tagName,
					"user_id": userID,
				}).Error("Error inserting new tag")
				if !isExternalTx {
					txInstance.Rollback()
				}
				return err
			}
		}
	}

	// Only commit if this function created the transaction
	if !isExternalTx {
		err = txInstance.Commit()
		if err != nil {
			logrs.WithError(err).Error("Error committing transaction for adding tags")
			return err
		}
	}

	return nil
}
