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
	maxTagNameLength = 255
)

type Tag struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginationParams struct {
	Limit  int
	Offset int
}

func InsertTag(ctx context.Context, tag *Tag, dbService *DBService, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(dbService, otx...)

	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > maxTagNameLength {
		return errors.New("invalid input for tag")
	}

	tag.ID, _ = util.GenerateSnowflakeID()
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	query, args, err := squirrel.Insert("tags").
		Columns("id", "user_id", "name", "created_at", "updated_at").
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

func UpdateTag(ctx context.Context, tag *Tag, dbService *DBService, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(dbService, otx...)

	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > maxTagNameLength {
		return errors.New("invalid input for tag")
	}

	tag.UpdatedAt = time.Now()

	query, args, err := squirrel.Update("tags").
		Set("name", tag.Name).
		Set("updated_at", tag.UpdatedAt).
		Where(squirrel.Eq{"id": tag.ID, "user_id": tag.UserID}).
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

func DeleteTag(ctx context.Context, tagID int64, userID int64, dbService *DBService, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(dbService, otx...)

	query, args, err := squirrel.Delete("tags").
		Where(squirrel.Eq{"id": tagID, "user_id": userID}).
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

func GetTagByID(ctx context.Context, tagID int64, userID int64, dbService *DBService, otx ...*sql.Tx) (*Tag, error) {
	_, executor := getExecutor(dbService, otx...)

	query, args, err := squirrel.Select("id", "user_id", "name", "created_at", "updated_at").
		From("tags").
		Where(squirrel.Eq{"id": tagID, "user_id": userID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query for retrieving tag by ID")
	}

	row := executor.QueryRowContext(ctx, query, args...)
	tag := &Tag{}
	err = row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, errors.Wrapf(err, "failed to retrieve tag by ID: %d", tagID)
	}

	return tag, nil
}

func GetAllTags(ctx context.Context, userID int64, pagination PaginationParams, dbService *DBService, otx ...*sql.Tx) ([]Tag, error) {
	_, executor := getExecutor(dbService, otx...)

	query, args, err := squirrel.Select("id", "user_id", "name", "created_at", "updated_at").
		From("tags").
		Where(squirrel.Eq{"user_id": userID}).
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

	var tags []Tag
	for rows.Next() {
		var tag Tag
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

func GetTagByName(ctx context.Context, name string, userID int64, dbService *DBService, otx ...*sql.Tx) (*Tag, error) {
	_, executor := getExecutor(dbService, otx...)

	query, args, err := squirrel.Select("id", "user_id", "name", "created_at", "updated_at").
		From("tags").
		Where(squirrel.Eq{"name": name, "user_id": userID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query for retrieving tag by name")
	}

	row := executor.QueryRowContext(ctx, query, args...)
	tag := &Tag{}
	err = row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tag not found")
		}
		return nil, errors.Wrapf(err, "failed to retrieve tag by name: %s for userID: %d", name, userID)
	}

	return tag, nil
}
