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
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type TagModel struct {
	TableTags        string
	ColumnID         string
	ColumnUserID     string
	ColumnName       string
	ColumnCreatedAt  string
	ColumnUpdatedAt  string
	MaxTagNameLength int
}

func NewTagModel() *TagModel {
	return &TagModel{
		TableTags:        "tags",
		ColumnID:         "id",
		ColumnUserID:     "user_id",
		ColumnName:       "name",
		ColumnCreatedAt:  "created_at",
		ColumnUpdatedAt:  "updated_at",
		MaxTagNameLength: 255, // Adjust as per your requirement
	}
}

func (tm *TagModel) InsertTag(ctx context.Context, tag *interfaces.Tag, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > tm.MaxTagNameLength {
		return errors.New("invalid input for tag")
	}

	tag.ID, _ = util.GenerateSnowflakeID()
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	query, args, err := squirrel.Insert(tm.TableTags).
		Columns(tm.ColumnID, tm.ColumnUserID, tm.ColumnName, tm.ColumnCreatedAt, tm.ColumnUpdatedAt).
		Values(tag.ID, tag.UserID, tag.Name, tag.CreatedAt, tag.UpdatedAt).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build insert query for tag")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to insert tag: %v", tag)
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (tm *TagModel) UpdateTag(ctx context.Context, tag *interfaces.Tag, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > tm.MaxTagNameLength {
		return errors.New("invalid input for tag")
	}

	tag.UpdatedAt = time.Now()

	query, args, err := squirrel.Update(tm.TableTags).
		Set(tm.ColumnName, tag.Name).
		Set(tm.ColumnUpdatedAt, tag.UpdatedAt).
		Where(squirrel.Eq{tm.ColumnID: tag.ID, tm.ColumnUserID: tag.UserID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build update query for tag")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to update tag: %v", tag)
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (tm *TagModel) DeleteTag(ctx context.Context, tagID int64, userID int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	query, args, err := squirrel.Delete(tm.TableTags).
		Where(squirrel.Eq{tm.ColumnID: tagID, tm.ColumnUserID: userID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build delete query for tag")
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to delete tag with tagID: %d and userID: %d", tagID, userID)
	}

	commitOrRollback(executor, isExternalTx, err)

	return nil
}

func (tm *TagModel) GetTagByID(ctx context.Context, tagID int64, userID int64, otx ...*sql.Tx) (*interfaces.Tag, error) {
	_, executor := getExecutor(otx...)

	query, args, err := squirrel.Select(tm.ColumnID, tm.ColumnUserID, tm.ColumnName, tm.ColumnCreatedAt, tm.ColumnUpdatedAt).
		From(tm.TableTags).
		Where(squirrel.Eq{tm.ColumnID: tagID, tm.ColumnUserID: userID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query for retrieving tag by ID")
	}

	row := executor.QueryRowContext(ctx, query, args...)
	tag := &interfaces.Tag{}
	err = row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, errors.Wrapf(err, "failed to retrieve tag by ID: %d", tagID)
	}

	return tag, nil
}

func (tm *TagModel) GetAllTags(ctx context.Context, userID int64, pagination interfaces.PaginationParams, otx ...*sql.Tx) ([]interfaces.Tag, error) {
	_, executor := getExecutor(otx...)

	query, args, err := squirrel.Select(tm.ColumnID, tm.ColumnUserID, tm.ColumnName, tm.ColumnCreatedAt, tm.ColumnUpdatedAt).
		From(tm.TableTags).
		Where(squirrel.Eq{tm.ColumnUserID: userID}).
		Limit(uint64(pagination.Limit)).
		Offset(uint64(pagination.Offset)).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query for retrieving all tags")
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve all tags for userID: %d", userID)
	}
	defer rows.Close()

	var tags []interfaces.Tag
	for rows.Next() {
		var tag interfaces.Tag
		err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan tag")
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate over all tags")
	}

	return tags, nil
}

func (tm *TagModel) GetTagByName(ctx context.Context, name string, userID int64, otx ...*sql.Tx) (*interfaces.Tag, error) {
	_, executor := getExecutor(otx...)

	query, args, err := squirrel.Select(tm.ColumnID, tm.ColumnUserID, tm.ColumnName, tm.ColumnCreatedAt, tm.ColumnUpdatedAt).
		From(tm.TableTags).
		Where(squirrel.Eq{tm.ColumnName: name, tm.ColumnUserID: userID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query for retrieving tag by name")
	}

	row := executor.QueryRowContext(ctx, query, args...)
	tag := &interfaces.Tag{}
	err = row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, errors.Wrapf(err, "failed to retrieve tag by name: %s for userID: %d", name, userID)
	}

	return tag, nil
}
