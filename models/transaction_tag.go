package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type TransactionTag struct {
	TransactionID int64     `json:"transaction_id"`
	TagID         int64     `json:"tag_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GetTagsByTransactionID(ctx context.Context, transactionID int64, otx ...*sql.Tx) ([]Tag, error) {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting transaction")
	}
	sql, args, err := squirrel.Select("t.id", "t.name").
		From("tags t").
		Join("transaction_tags tt ON t.id = tt.tag_id").
		Where(squirrel.Eq{"tt.transaction_id": transactionID}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build SQL query for GetTagsByTransactionID")
	}

	rows, err := tx.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error querying tags for transaction")
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, errors.Wrap(err, "error scanning tag row")
		}
		tags = append(tags, tag)
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return nil, errors.Wrap(err, "committing transaction failed")
		}
	}
	return tags, nil
}

func InsertTransactionTag(ctx context.Context, transactionID, tagID int64, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}
	sql, args, err := squirrel.Insert("transaction_tags").
		Columns("transaction_id", "tag_id", "created_at", "updated_at").
		Values(transactionID, tagID, time.Now(), time.Now()).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build SQL query for InsertTransactionTag")
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "error inserting transaction tag")
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}

func DeleteTransactionTag(ctx context.Context, transactionID, tagID int64, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	sql, args, err := squirrel.Delete("transaction_tags").
		Where(squirrel.Eq{"transaction_id": transactionID, "tag_id": tagID}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build SQL query for DeleteTransactionTag")
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "error deleting transaction tag")
	}

	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}

	return nil
}

func AddTagsToTransaction(ctx context.Context, transactionID int64, tags []string, userID int64, otx ...*sql.Tx) error {

	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	for _, tagName := range tags {
		tag, err := GetTagByName(ctx, tagName, userID, tx)
		if err != nil {
			return errors.Wrap(err, "error getting tag by name")
		}
		err = InsertTransactionTag(ctx, transactionID, tag.ID, tx)
		if err != nil {
			return errors.Wrap(err, "error associating tag with transaction")
		}
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}

func UpdateTagsForTransaction(ctx context.Context, transactionID int64, tags []string, userID int64, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	err = DeleteTagsFromTransaction(ctx, transactionID, tx)
	if err != nil {
		return errors.Wrap(err, "error removing existing tags from transaction")
	}

	err = AddTagsToTransaction(ctx, transactionID, tags, userID, tx)
	if err != nil {
		return errors.Wrap(err, "error adding new tags to transaction")
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}

func DeleteTagsFromTransaction(ctx context.Context, transactionID int64, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}
	sql, args, err := squirrel.Delete("transaction_tags").
		Where(squirrel.Eq{"transaction_id": transactionID}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build SQL query for DeleteTagsFromTransaction")
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "error deleting tags from transaction")
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}
