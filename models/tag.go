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

func InsertTag(ctx context.Context, tag *Tag, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > maxTagNameLength {
		return errors.New("invalid input for tag")
	}

	tag.ID, _ = util.GenerateSnowflakeID()
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	sql, args, err := squirrel.Insert("tags").
		Columns("id", "user_id", "name", "created_at", "updated_at").
		Values(tag.ID, tag.UserID, tag.Name, tag.CreatedAt, tag.UpdatedAt).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build insert query for tag")
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to insert tag: %v", tag)
	}

	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}

func UpdateTag(ctx context.Context, tag *Tag, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > maxTagNameLength {
		return errors.New("invalid input for tag")
	}

	tag.UpdatedAt = time.Now()

	sql, args, err := squirrel.Update("tags").
		Set("name", tag.Name).
		Set("updated_at", tag.UpdatedAt).
		Where(squirrel.Eq{"id": tag.ID, "user_id": tag.UserID}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build update query for tag")
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to update tag: %v", tag)
	}

	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}

	return nil
}

func DeleteTag(ctx context.Context, tagID int64, userID int64, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	sql, args, err := squirrel.Delete("tags").
		Where(squirrel.Eq{"id": tagID, "user_id": userID}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "failed to build delete query for tag")
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to delete tag with tagID: %d and userID: %d", tagID, userID)
	}

	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}

	return nil
}

func GetTagByID(ctx context.Context, tagID int64, userID int64, otx ...*sql.Tx) (*Tag, error) {

	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting transaction")
	}

	sql, args, err := squirrel.Select("id", "user_id", "name", "created_at", "updated_at").
		From("tags").
		Where(squirrel.Eq{"id": tagID, "user_id": userID}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query for retrieving tag by ID")
	}

	row := tx.QueryRowContext(ctx, sql, args...)
	tag := &Tag{}
	err = row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve tag by ID: %d", tagID)
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return nil, errors.Wrap(err, "committing transaction failed")
		}
	}

	return tag, nil
}

func GetAllTags(ctx context.Context, userID int64, pagination PaginationParams, otx ...*sql.Tx) ([]Tag, error) {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting transaction")
	}
	sql, args, err := squirrel.Select("id", "user_id", "name", "created_at", "updated_at").
		From("tags").
		Where(squirrel.Eq{"user_id": userID}).
		Limit(uint64(pagination.Limit)).
		Offset(uint64(pagination.Offset)).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query for retrieving all tags")
	}

	rows, err := tx.QueryContext(ctx, sql, args...)
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

	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return nil, errors.Wrap(err, "committing transaction failed")
		}
	}
	return tags, nil
}

func GetTagByName(ctx context.Context, name string, userID int64, otx ...*sql.Tx) (*Tag, error) {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting transaction")
	}
	sql, args, err := squirrel.Select("id", "user_id", "name", "created_at", "updated_at").
		From("tags").
		Where(squirrel.Eq{"name": name, "user_id": userID}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "failed to build query for retrieving tag by name")
	}

	row := tx.QueryRowContext(ctx, sql, args...)
	tag := &Tag{}
	err = row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve tag by name: %s for userID: %d", name, userID)
	}

	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return nil, errors.Wrap(err, "committing transaction failed")
		}
	}
	return tag, nil
}

func executeTxQuery(ctx context.Context, tx *sql.Tx, sqlStr string, args ...interface{}) error {
	_, err := tx.ExecContext(ctx, sqlStr, args...)
	return err
}

func executeQuery(ctx context.Context, db *sql.DB, sqlStr string, args ...interface{}) error {
	_, err := db.ExecContext(ctx, sqlStr, args...)
	return err
}

func executeQueryRow(ctx context.Context, db *sql.DB, sqlStr string, args ...interface{}) (*Tag, error) {
	row := db.QueryRowContext(ctx, sqlStr, args...)
	tag := &Tag{}
	err := row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func executeQueryRows(ctx context.Context, db *sql.DB, sqlStr string, args ...interface{}) ([]Tag, error) {
	rows, err := db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func executeTxQueryRow(ctx context.Context, tx *sql.Tx, sqlStr string, args ...interface{}) (*Tag, error) {
	row := tx.QueryRowContext(ctx, sqlStr, args...)
	tag := &Tag{}
	err := row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return tag, nil
}
