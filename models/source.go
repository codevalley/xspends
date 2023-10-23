package models

import (
	"errors"
	"log"
)

const (
	SourceTypeCredit  = "CREDIT"
	SourceTypeSavings = "SAVINGS"
)

type Source struct {
	ID      string  `json:"id"`
	UserID  string  `json:"user_id"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}

func InsertSource(source Source) error {
	// Validate fields
	if source.Name == "" || source.UserID == "" {
		return errors.New("source name and user ID cannot be empty")
	}

	if source.Type != SourceTypeCredit && source.Type != SourceTypeSavings {
		return errors.New("invalid source type; only 'CREDIT' and 'SAVINGS' are allowed")
	}

	_, err := GetDB().Exec("INSERT INTO sources (id, user_id, name, type, balance) VALUES (?, ?, ?, ?, ?)", source.ID, source.UserID, source.Name, source.Type, source.Balance)
	if err != nil {
		log.Println("Error inserting source:", err)
		return err
	}
	return nil
}

func UpdateSource(source Source) error {
	// Validate fields
	if source.Name == "" || source.UserID == "" {
		return errors.New("source name and user ID cannot be empty")
	}

	if source.Type != SourceTypeCredit && source.Type != SourceTypeSavings {
		return errors.New("invalid source type; only 'CREDIT' and 'SAVINGS' are allowed")
	}

	_, err := GetDB().Exec("UPDATE sources SET user_id=?, name=?, type=?, balance=? WHERE id=?", source.UserID, source.Name, source.Type, source.Balance, source.ID)
	if err != nil {
		log.Println("Error updating source:", err)
		return err
	}
	return nil
}

func DeleteSource(sourceID string) error {
	_, err := GetDB().Exec("DELETE FROM sources WHERE id=?", sourceID)
	if err != nil {
		log.Println("Error deleting source:", err)
		return err
	}
	return nil
}

func GetSourceByID(sourceID string) (*Source, error) {
	row := GetDB().QueryRow("SELECT id, user_id, name, type, balance FROM sources WHERE id=?", sourceID)
	var source Source
	err := row.Scan(&source.ID, &source.UserID, &source.Name, &source.Type, &source.Balance)
	if err != nil {
		log.Println("Error retrieving source by ID:", err)
		return nil, err
	}
	return &source, nil
}
