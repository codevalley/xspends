package models

import (
	"database/sql"
	"time"

	"xspends/util"

	"github.com/Masterminds/squirrel"
)

const (
	SourceTypeCredit  = "CREDIT"
	SourceTypeSavings = "SAVINGS"
)

type Source struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func InsertSource(source *Source) error {
	if source.Name == "" || source.UserID == 0 {
		return util.ErrInvalidInput
	}
	if source.Type != SourceTypeCredit && source.Type != SourceTypeSavings {
		return util.ErrInvalidType
	}

	source.ID, _ = util.GenerateSnowflakeID()
	source.CreatedAt = time.Now()
	source.UpdatedAt = source.CreatedAt

	query, args, err := SQLBuilder.Insert("sources").
		Columns("id", "user_id", "name", "type", "balance", "created_at", "updated_at").
		Values(source.ID, source.UserID, source.Name, source.Type, source.Balance, source.CreatedAt, source.UpdatedAt).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Error preparing insert SQL for source")
		return util.ErrDatabase
	}

	_, err = GetDB().Exec(query, args...)
	if err != nil {
		logrs.WithError(err).Error("Error executing insert for source")
		return util.ErrDatabase
	}

	return nil
}

func UpdateSource(source *Source) error {
	if source.Name == "" || source.UserID == 0 {
		return util.ErrInvalidInput
	}
	if source.Type != SourceTypeCredit && source.Type != SourceTypeSavings {
		return util.ErrInvalidType
	}

	source.UpdatedAt = time.Now()

	query, args, err := SQLBuilder.Update("sources").
		Set("name", source.Name).
		Set("type", source.Type).
		Set("balance", source.Balance).
		Set("updated_at", source.UpdatedAt).
		Where(squirrel.Eq{"id": source.ID, "user_id": source.UserID}).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Error preparing update SQL for source")
		return util.ErrDatabase
	}

	_, err = GetDB().Exec(query, args...)
	if err != nil {
		logrs.WithError(err).Error("Error executing update for source")
		return util.ErrDatabase
	}

	return nil
}

func DeleteSource(sourceID int64, userID int64) error {
	query, args, err := SQLBuilder.Delete("sources").
		Where(squirrel.Eq{"id": sourceID, "user_id": userID}).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Error preparing delete SQL for source")
		return util.ErrDatabase
	}

	_, err = GetDB().Exec(query, args...)
	if err != nil {
		logrs.WithError(err).Error("Error executing delete for source")
		return util.ErrDatabase
	}

	return nil
}

func GetSourceByID(sourceID int64, userID int64) (*Source, error) {
	query, args, err := SQLBuilder.Select("id", "user_id", "name", "type", "balance", "created_at", "updated_at").
		From("sources").
		Where(squirrel.Eq{"id": sourceID, "user_id": userID}).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Error preparing select SQL for source by ID")
		return nil, util.ErrDatabase
	}

	source := &Source{}
	err = GetDB().QueryRow(query, args...).Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance, &source.CreatedAt, &source.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, util.ErrSourceNotFound
		}
		logrs.WithError(err).Error("Error querying source by ID")
		return nil, util.ErrDatabase
	}

	return source, nil
}

func GetSources(userID int64) ([]Source, error) {
	query, args, err := SQLBuilder.Select("id", "user_id", "name", "type", "balance", "created_at", "updated_at").
		From("sources").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Error preparing select SQL for sources by user ID")
		return nil, util.ErrDatabase
	}

	rows, err := GetDB().Query(query, args...)
	if err != nil {
		logrs.WithError(err).Error("Error querying sources by user ID")
		return nil, util.ErrDatabase
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var source Source
		if err = rows.Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance, &source.CreatedAt, &source.UpdatedAt); err != nil {
			logrs.WithError(err).Error("Error scanning source row")
			return nil, util.ErrDatabase
		}
		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		logrs.WithError(err).Error("Error during row processing for sources")
		return nil, util.ErrDatabase
	}

	return sources, nil
}

func SourceIDExists(sourceID int64, userID int64) (bool, error) {
	query, args, err := SQLBuilder.Select("1").
		From("sources").
		Where(squirrel.Eq{"id": sourceID, "user_id": userID}).
		Limit(1).
		ToSql()

	if err != nil {
		logrs.WithError(err).Error("Error preparing SQL to check if source exists by ID")
		return false, util.ErrDatabase
	}

	var exists int
	err = GetDB().QueryRow(query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logrs.WithError(err).Error("Error checking if source exists by ID")
		return false, util.ErrDatabase
	}

	return exists == 1, nil
}
