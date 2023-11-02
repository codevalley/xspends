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

func InsertSource(ctx context.Context, source *Source) error {
	if source.Name == "" || source.UserID == 0 {
		return errors.New("invalid input: name or user ID is empty")
	}
	if source.Type != SourceTypeCredit && source.Type != SourceTypeSavings {
		return errors.New("invalid type: type must be CREDIT or SAVINGS")
	}

	source.ID, _ = util.GenerateSnowflakeID()
	source.CreatedAt = time.Now()
	source.UpdatedAt = source.CreatedAt

	query, args, err := SQLBuilder.Insert("sources").
		Columns("id", "user_id", "name", "type", "balance", "created_at", "updated_at").
		Values(source.ID, source.UserID, source.Name, source.Type, source.Balance, source.CreatedAt, source.UpdatedAt).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "preparing insert SQL for source")
	}

	_, err = GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing insert for source")
	}

	return nil
}

func UpdateSource(ctx context.Context, source *Source) error {
	if source.Name == "" || source.UserID == 0 {
		return errors.New("invalid input: name or user ID is empty")
	}
	if source.Type != SourceTypeCredit && source.Type != SourceTypeSavings {
		return errors.New("invalid type: type must be CREDIT or SAVINGS")
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
		return errors.Wrap(err, "preparing update SQL for source")
	}

	_, err = GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing update for source")
	}

	return nil
}

func DeleteSource(ctx context.Context, sourceID int64, userID int64) error {
	query, args, err := SQLBuilder.Delete("sources").
		Where(squirrel.Eq{"id": sourceID, "user_id": userID}).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "preparing delete SQL for source")
	}

	_, err = GetDB().ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "executing delete for source")
	}

	return nil
}

func GetSourceByID(ctx context.Context, sourceID int64, userID int64) (*Source, error) {
	query, args, err := SQLBuilder.Select("id", "user_id", "name", "type", "balance", "created_at", "updated_at").
		From("sources").
		Where(squirrel.Eq{"id": sourceID, "user_id": userID}).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "preparing select SQL for source by ID")
	}

	source := &Source{}
	err = GetDB().QueryRowContext(ctx, query, args...).Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance, &source.CreatedAt, &source.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("source not found")
		}
		return nil, errors.Wrap(err, "querying source by ID")
	}

	return source, nil
}

func GetSources(ctx context.Context, userID int64) ([]Source, error) {
	query, args, err := SQLBuilder.Select("id", "user_id", "name", "type", "balance", "created_at", "updated_at").
		From("sources").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "preparing select SQL for sources by user ID")
	}

	rows, err := GetDB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying sources by user ID")
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var source Source
		if err = rows.Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance, &source.CreatedAt, &source.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "scanning source row")
		}
		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "during row processing for sources")
	}

	return sources, nil
}

func SourceIDExists(ctx context.Context, sourceID int64, userID int64) (bool, error) {
	query, args, err := SQLBuilder.Select("1").
		From("sources").
		Where(squirrel.Eq{"id": sourceID, "user_id": userID}).
		Limit(1).
		ToSql()

	if err != nil {
		return false, errors.Wrap(err, "preparing SQL to check if source exists by ID")
	}

	var exists int
	err = GetDB().QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "checking if source exists by ID")
	}

	return exists == 1, nil
}
