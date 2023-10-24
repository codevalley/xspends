package models

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"xspends/util"
)

const (
	SourceTypeCredit  = "CREDIT"
	SourceTypeSavings = "SAVINGS"
)

var (
	ErrSourceNotFound = errors.New("source not found")
	ErrInvalidType    = errors.New("invalid source type; only 'CREDIT' and 'SAVINGS' are allowed")
)

type Source struct {
	ID        int64     `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func InsertSource(source Source) error {
	// Validation
	if source.Name == "" || source.UserID == 0 {
		return errors.New("source name and user ID cannot be empty")
	}
	if source.Type != SourceTypeCredit && source.Type != SourceTypeSavings {
		return ErrInvalidType
	}

	sid, err := util.GenerateSnowflakeID()
	if err != nil {
		log.Printf("[ERROR] Generating snowflake ID for source: %v", err)
		return err
	}
	source.ID = sid
	source.CreatedAt = time.Now()
	source.UpdatedAt = time.Now()

	stmt, err := GetDB().Prepare("INSERT INTO sources (id, user_id, name, type, balance, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[ERROR] Preparing insert statement for source: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(source.ID, source.UserID, source.Name, source.Type, source.Balance, source.CreatedAt, source.UpdatedAt)
	if err != nil {
		log.Printf("[ERROR] Inserting source: %v", err)
		return err
	}

	return nil
}

func UpdateSource(source Source) error {
	// Validation
	if source.Name == "" || source.UserID == 0 {
		return errors.New("source name and user ID cannot be empty")
	}
	if source.Type != SourceTypeCredit && source.Type != SourceTypeSavings {
		return ErrInvalidType
	}

	source.UpdatedAt = time.Now()

	stmt, err := GetDB().Prepare("UPDATE sources SET user_id=?, name=?, type=?, balance=?, updated_at=? WHERE id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing update statement for source: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(source.UserID, source.Name, source.Type, source.Balance, source.UpdatedAt, source.ID)
	if err != nil {
		log.Printf("[ERROR] Updating source: %v", err)
		return err
	}

	return nil
}

func DeleteSource(sourceID int64) error {
	stmt, err := GetDB().Prepare("DELETE FROM sources WHERE id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing delete statement for source: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sourceID)
	if err != nil {
		log.Printf("[ERROR] Deleting source: %v", err)
		return err
	}

	return nil
}

func GetSourceByID(sourceID int64) (*Source, error) {
	stmt, err := GetDB().Prepare("SELECT id, user_id, name, type, balance, created_at, updated_at FROM sources WHERE id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing select statement for source by ID: %v", err)
		return nil, err
	}
	defer stmt.Close()

	source := &Source{}
	err = stmt.QueryRow(sourceID).Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance, &source.CreatedAt, &source.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSourceNotFound
		}
		log.Printf("[ERROR] Retrieving source by ID: %v", err)
		return nil, err
	}

	return source, nil
}

func GetSourcesByUserID(userID int) ([]Source, error) {
	stmt, err := GetDB().Prepare("SELECT id, user_id, name, type, balance, created_at, updated_at FROM sources WHERE user_id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing select statement for sources by user ID: %v", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		log.Printf("[ERROR] Querying sources by user ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var source Source
		if err := rows.Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance, &source.CreatedAt, &source.UpdatedAt); err != nil {
			log.Printf("[ERROR] Scanning source row: %v", err)
			return nil, err
		}
		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[ERROR] During rows scan: %v", err)
		return nil, err
	}

	return sources, nil
}
