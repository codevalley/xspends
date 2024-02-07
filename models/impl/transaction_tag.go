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

package impl

import (
	"context"
	"database/sql"
	"time"
	"xspends/models/interfaces"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type TransactionTagModel struct {
	TableTransactionTags string
	ColumnTransactionID  string
	ColumnTagID          string
	ColumnCreatedAt      string
	ColumnUpdatedAt      string
}

func NewTransactionTagModel() *TransactionTagModel {
	return &TransactionTagModel{
		TableTransactionTags: "transaction_tags",
		ColumnTransactionID:  "transaction_id",
		ColumnTagID:          "tag_id",
		ColumnCreatedAt:      "created_at",
		ColumnUpdatedAt:      "updated_at",
	}
}

func (tm *TransactionTagModel) GetTagsByTransactionID(ctx context.Context, transactionID int64, otx ...*sql.Tx) ([]interfaces.Tag, error) {
	_, executor := getExecutor(otx...)

	//TODO: Has to resolved, can't be hardcoded like this.
	tagID := "tag_id"
	tagName := "name"
	//TODO: possible logical bug here. Review and close
	query, args, err := squirrel.Select(tagID, tagName).
		From("tags t").
		Join(tm.TableTransactionTags + " tt ON t." + tagID + " = tt." + tm.ColumnTagID).
		Where(squirrel.Eq{"tt." + tm.ColumnTransactionID: transactionID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build SQL query for GetTagsByTransactionID")
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error querying tags for transaction")
	}
	defer rows.Close()

	var tags []interfaces.Tag
	for rows.Next() {
		var tag interfaces.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, errors.Wrap(err, "error scanning tag row")
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (tm *TransactionTagModel) InsertTransactionTag(ctx context.Context, transactionID, tagID int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	query, args, err := squirrel.Insert(tm.TableTransactionTags).
		Columns(tm.ColumnTransactionID, tm.ColumnTagID, tm.ColumnCreatedAt, tm.ColumnUpdatedAt).
		Values(transactionID, tagID, time.Now(), time.Now()).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build SQL query for InsertTransactionTag")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "error inserting transaction tag")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (tm *TransactionTagModel) DeleteTransactionTag(ctx context.Context, transactionID, tagID int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	query, args, err := squirrel.Delete(tm.TableTransactionTags).
		Where(squirrel.Eq{tm.ColumnTransactionID: transactionID, tm.ColumnTagID: tagID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build SQL query for DeleteTransactionTag")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "error deleting transaction tag")
	}

	commitOrRollback(executor, isExternalTx, err)

	return nil
}

func (tm *TransactionTagModel) AddTagsToTransaction(ctx context.Context, transactionID int64, tags []string, scopes []int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	for _, tagName := range tags {
		tag, err := GetModelsService().TagModel.GetTagByName(ctx, tagName, scopes, otx...)
		if err != nil {
			return errors.Wrap(err, "error getting tag by name")
		}
		err = GetModelsService().TransactionTagModel.InsertTransactionTag(ctx, transactionID, tag.ID, otx...)
		if err != nil {
			return errors.Wrap(err, "error associating tag with transaction")
		}
	}

	commitOrRollback(executor, isExternalTx, nil)
	return nil
}

func (tm *TransactionTagModel) UpdateTagsForTransaction(ctx context.Context, transactionID int64, tags []string, scopes []int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	err := GetModelsService().TransactionTagModel.DeleteTagsFromTransaction(ctx, transactionID, otx...)
	if err != nil {
		return errors.Wrap(err, "error removing existing tags from transaction")
	}

	err = GetModelsService().TransactionTagModel.AddTagsToTransaction(ctx, transactionID, tags, scopes, otx...)
	if err != nil {
		return errors.Wrap(err, "error adding new tags to transaction")
	}

	commitOrRollback(executor, isExternalTx, nil)
	return nil
}

func (tm *TransactionTagModel) DeleteTagsFromTransaction(ctx context.Context, transactionID int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	query, args, err := squirrel.Delete(tm.TableTransactionTags).
		Where(squirrel.Eq{tm.ColumnTransactionID: transactionID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build SQL query for DeleteTagsFromTransaction")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "error deleting tags from transaction")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}
