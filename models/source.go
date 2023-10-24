package models

import (
	"database/sql"
	"errors"
	"log"
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
	ID      int64   `json:"id"`
	UserID  int     `json:"user_id"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
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
		log.Printf("Error generating user: %v", err)
		return err
	}
	source.ID = sid

	stmt, err := GetDB().Prepare("INSERT INTO sources (id, user_id, name, type, balance) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing insert statement for source:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(source.ID, source.UserID, source.Name, source.Type, source.Balance)
	if err != nil {
		log.Println("Error inserting source:", err)
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

	stmt, err := GetDB().Prepare("UPDATE sources SET user_id=?, name=?, type=?, balance=? WHERE id=?")
	if err != nil {
		log.Println("Error preparing update statement for source:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(source.UserID, source.Name, source.Type, source.Balance, source.ID)
	if err != nil {
		log.Println("Error updating source:", err)
		return err
	}

	return nil
}

func DeleteSource(sourceID int) error {
	stmt, err := GetDB().Prepare("DELETE FROM sources WHERE id=?")
	if err != nil {
		log.Println("Error preparing delete statement for source:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sourceID)
	if err != nil {
		log.Println("Error deleting source:", err)
		return err
	}

	return nil
}

func GetSourceByID(sourceID int) (*Source, error) {
	stmt, err := GetDB().Prepare("SELECT id, user_id, name, type, balance FROM sources WHERE id=?")
	if err != nil {
		log.Println("Error preparing select statement for source by ID:", err)
		return nil, err
	}
	defer stmt.Close()

	source := &Source{}
	err = stmt.QueryRow(sourceID).Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSourceNotFound
		}
		log.Println("Error retrieving source by ID:", err)
		return nil, err
	}

	return source, nil
}

func GetSourcesByUserID(userID int) ([]Source, error) {
	stmt, err := GetDB().Prepare("SELECT id, user_id, name, type, balance FROM sources WHERE user_id=?")
	if err != nil {
		log.Println("Error preparing select statement for sources by user ID:", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		log.Println("Error querying sources by user ID:", err)
		return nil, err
	}
	defer rows.Close()

	var sources []Source
	for rows.Next() {
		var source Source
		if err := rows.Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance); err != nil {
			log.Println("Error scanning source row:", err)
			return nil, err
		}
		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error during rows scan:", err)
		return nil, err
	}

	return sources, nil
}
