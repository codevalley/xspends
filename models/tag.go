package models

import (
	"database/sql"
	"time"
	"xspends/util" // Adjust this import to your project's structure

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
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

// PaginationParams holds parameters for paginating database queries
type PaginationParams struct {
	Limit  int
	Offset int
}

// InsertTag adds a new tag to the database.
func InsertTag(tag *Tag, tx ...*sql.Tx) error {
	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > maxTagNameLength {
		return util.ErrInvalidInput
	}

	tag.ID, _ = util.GenerateSnowflakeID()
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()

	sql, args, err := SQLBuilder.Insert("tags").
		Columns("id", "user_id", "name", "created_at", "updated_at").
		Values(tag.ID, tag.UserID, tag.Name, tag.CreatedAt, tag.UpdatedAt).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Failed to build insert query for tag")
		return util.ErrDatabase
	}

	err = executeTxQuery(tx, sql, args...)
	if err != nil {
		logrs.WithError(err).WithField("tag", tag).Error("Failed to insert tag")
		return util.ErrDatabase
	}

	return nil
}

// UpdateTag updates an existing tag in the database.
func UpdateTag(tag *Tag, tx ...*sql.Tx) error {
	if tag.UserID <= 0 || len(tag.Name) == 0 || len(tag.Name) > maxTagNameLength {
		return util.ErrInvalidInput
	}

	tag.UpdatedAt = time.Now()

	sql, args, err := SQLBuilder.Update("tags").
		Set("name", tag.Name).
		Set("updated_at", tag.UpdatedAt).
		Where(squirrel.Eq{"id": tag.ID, "user_id": tag.UserID}).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Failed to build update query for tag")
		return util.ErrDatabase
	}

	err = executeTxQuery(tx, sql, args...)
	if err != nil {
		logrs.WithError(err).WithField("tag", tag).Error("Failed to update tag")
		return util.ErrDatabase
	}

	return nil
}

// DeleteTag removes a tag from the database.
func DeleteTag(tagID int64, userID int64) error {
	sql, args, err := SQLBuilder.Delete("tags").
		Where(squirrel.Eq{"id": tagID, "user_id": userID}).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Failed to build delete query for tag")
		return util.ErrDatabase
	}

	err = executeQuery(sql, args...)
	if err != nil {
		logrs.WithError(err).WithFields(logrus.Fields{
			"tagID":  tagID,
			"userID": userID,
		}).Error("Failed to delete tag")
		return util.ErrDatabase
	}

	return nil
}

// GetTagByID retrieves a tag by its ID.
func GetTagByID(tagID int64, userID int64) (*Tag, error) {
	sql, args, err := SQLBuilder.Select("id", "user_id", "name", "created_at", "updated_at").
		From("tags").
		Where(squirrel.Eq{"id": tagID, "user_id": userID}).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Failed to build query for retrieving tag by ID")
		return nil, util.ErrDatabase
	}

	tag := &Tag{}
	err = executeQueryRow(sql, tag, args...)
	if err != nil {
		logrs.WithError(err).WithField("tagID", tagID).Error("Failed to retrieve tag by ID")
		return nil, err
	}

	return tag, nil
}

// GetAllTags retrieves all tags for a user with pagination.
func GetAllTags(userID int64, pagination PaginationParams) ([]Tag, error) {
	sql, args, err := SQLBuilder.Select("id", "user_id", "name", "created_at", "updated_at").
		From("tags").
		Where(squirrel.Eq{"user_id": userID}).
		Limit(uint64(pagination.Limit)).
		Offset(uint64(pagination.Offset)).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Failed to build query for retrieving all tags")
		return nil, util.ErrDatabase
	}

	tags := []Tag{}
	err = executeQueryRows(sql, &tags, args...)
	if err != nil {
		logrs.WithError(err).WithField("userID", userID).Error("Failed to retrieve all tags")
		return nil, err
	}

	return tags, nil
}

// GetTagByName retrieves a tag by its name for a specific user.
func GetTagByName(name string, userID int64, tx ...*sql.Tx) (*Tag, error) {
	sql, args, err := SQLBuilder.Select("id", "user_id", "name", "created_at", "updated_at").
		From("tags").
		Where(squirrel.Eq{"name": name, "user_id": userID}).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Failed to build query for retrieving tag by name")
		return nil, util.ErrDatabase
	}

	tag := &Tag{}
	err = executeTxQueryRow(tx, sql, tag, args...)
	if err != nil {
		logrs.WithError(err).WithFields(logrus.Fields{
			"name":   name,
			"userID": userID,
		}).Error("Failed to retrieve tag by name")
		return nil, err
	}

	return tag, nil
}

// Helper functions for executing SQL queries with or without transactions.
func executeTxQuery(tx []*sql.Tx, sqlStr string, args ...interface{}) error {
	var err error
	var res sql.Result

	if len(tx) > 0 {
		res, err = tx[0].Exec(sqlStr, args...)
	} else {
		db := GetDB()
		res, err = db.Exec(sqlStr, args...)
	}

	if err != nil {
		return err
	}

	_, err = res.LastInsertId() // Optionally use the result
	return err
}

func executeQuery(sql string, args ...interface{}) error {
	db := GetDB()
	res, err := db.Exec(sql, args...)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected() // Optionally use the result
	return err
}

func executeQueryRow(sql string, tag *Tag, args ...interface{}) error {
	db := GetDB()
	return db.QueryRow(sql, args...).Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
}

func executeQueryRows(sql string, tags *[]Tag, args ...interface{}) error {
	db := GetDB()
	rows, err := db.Query(sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return err
		}
		*tags = append(*tags, tag)
	}
	return rows.Err()
}

func executeTxQueryRow(tx []*sql.Tx, sqlStr string, tag *Tag, args ...interface{}) error {
	var row *sql.Row // This line needs to be changed.

	if len(tx) > 0 {
		// No need to declare a new variable. We can use := to get the row directly.
		row = tx[0].QueryRow(sqlStr, args...)
	} else {
		db := GetDB()
		// Same here, use := to get the row directly.
		row = db.QueryRow(sqlStr, args...)
	}

	// row is already the correct type, so just call Scan on it.
	err := row.Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	return err
}
