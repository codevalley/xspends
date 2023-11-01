package models

import (
	"database/sql"
	"time"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type TransactionTag struct {
	TransactionID int64     `json:"transaction_id"`
	TagID         int64     `json:"tag_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetTagsByTransactionID retrieves all tags for a specific transaction.
func GetTagsByTransactionID(transactionID int64) ([]Tag, error) {
	queryBuilder := squirrel.Select("t.id", "t.name").
		From("tags t").
		Join("transaction_tags tt ON t.id = tt.tag_id").
		Where(squirrel.Eq{"tt.transaction_id": transactionID})

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"error":         err,
		}).Error("Failed to build SQL query for GetTagsByTransactionID")
		return nil, util.ErrDatabase
	}

	rows, err := GetDB().Query(sql, args...)
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"error":         err,
		}).Error("Error querying tags for transaction")
		return nil, util.ErrDatabase
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			logrs.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error scanning tag row")
			return nil, util.ErrDatabase
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// InsertTransactionTag adds a new tag to a specific transaction.
func InsertTransactionTag(transactionID, tagID int64, tx ...*sql.Tx) error {
	queryBuilder := squirrel.Insert("transaction_tags").
		Columns("transaction_id", "tag_id", "created_at", "updated_at").
		Values(transactionID, tagID, time.Now(), time.Now())

	isExternalTx, txInstance, err := GetTransaction(tx...)
	if err != nil {
		return err
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"tagID":         tagID,
			"error":         err,
		}).Error("Failed to build SQL query for InsertTransactionTag")
		return util.ErrDatabase
	}

	_, err = txInstance.Exec(sql, args...)
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"tagID":         tagID,
			"error":         err,
		}).Error("Error inserting transaction tag")
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

// DeleteTransactionTag removes a specific tag from a specific transaction.
func DeleteTransactionTag(transactionID, tagID int64, tx ...*sql.Tx) error {
	queryBuilder := squirrel.Delete("transaction_tags").
		Where(squirrel.Eq{"transaction_id": transactionID, "tag_id": tagID})

	isExternalTx, txInstance, err := GetTransaction(tx...)
	if err != nil {
		return err
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"tagID":         tagID,
			"error":         err,
		}).Error("Failed to build SQL query for DeleteTransactionTag")
		return util.ErrDatabase
	}

	_, err = txInstance.Exec(sql, args...)
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"tagID":         tagID,
			"error":         err,
		}).Error("Error deleting transaction tag")
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

// AddTagsToTransaction adds multiple tags to a specific transaction.
func AddTagsToTransaction(transactionID int64, tags []string, userID int64, tx ...*sql.Tx) error {
	isExternalTx, txInstance, err := GetTransaction(tx...)
	if err != nil {
		return err
	}

	for _, tagName := range tags {
		tag, err := GetTagByName(tagName, userID, txInstance)
		if err != nil {
			logrs.WithFields(logrus.Fields{
				"tagName": tagName,
				"userID":  userID,
				"error":   err,
			}).Error("Error getting tag by name")
			if !isExternalTx {
				txInstance.Rollback()
			}
			return err
		}
		err = InsertTransactionTag(transactionID, tag.ID, txInstance)
		if err != nil {
			logrs.WithFields(logrus.Fields{
				"tagName":       tagName,
				"transactionID": transactionID,
				"error":         err,
			}).Error("Error associating tag with transaction")
			if !isExternalTx {
				txInstance.Rollback()
			}
			return err
		}
	}

	if !isExternalTx {
		return txInstance.Commit()
	}
	return nil
}

// UpdateTagsForTransaction updates the tag associations for a specific transaction.
func UpdateTagsForTransaction(transactionID int64, tags []string, userID int64, tx ...*sql.Tx) error {
	isExternalTx, txInstance, err := GetTransaction(tx...)
	if err != nil {
		return err
	}

	// Delete existing tags
	err = DeleteTagsFromTransaction(transactionID, txInstance)
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"error":         err,
		}).Error("Error removing existing tags from transaction")
		if !isExternalTx {
			txInstance.Rollback()
		}
		return err
	}

	// Add new tags
	err = AddTagsToTransaction(transactionID, tags, userID, txInstance)
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"tags":          tags,
			"error":         err,
		}).Error("Error adding new tags to transaction")
		if !isExternalTx {
			txInstance.Rollback()
		}
		return err
	}

	if !isExternalTx {
		return txInstance.Commit()
	}
	return nil
}

// RemoveTagsFromTransaction removes all tag associations from a specific transaction.
func DeleteTagsFromTransaction(transactionID int64, tx *sql.Tx) error {
	queryBuilder := squirrel.Delete("transaction_tags").
		Where(squirrel.Eq{"transaction_id": transactionID})

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"error":         err,
		}).Error("Failed to build SQL query for DeleteTagsFromTransaction")
		return util.ErrDatabase
	}

	_, err = tx.Exec(sql, args...)
	if err != nil {
		logrs.WithFields(logrus.Fields{
			"transactionID": transactionID,
			"error":         err,
		}).Error("Error deleting tags from transaction")
		return util.ErrDatabase
	}

	return nil
}
