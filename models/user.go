package models

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"xspends/util"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Currency string `json:"currency"`
	Password string `json:"-"`
}

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email already exists")
	ErrUsernameTaken = errors.New("username already exists")
)

func InsertUser(user *User) error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("mandatory fields missing")
	}

	// Generate SnowflakeId for user ID
	sid, err := util.GenerateSnowflakeID()
	if err != nil {
		log.Printf("Error generating user: %v", err)
		return err
	}
	user.ID = sid
	stmt, err := GetDB().Prepare("INSERT INTO users (id, username, name, email, currency, password) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Username, user.Name, user.Email, user.Currency, user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "username") {
			log.Printf("Username taken: %v", err)
			return ErrUsernameTaken
		}
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "email") {
			log.Printf("Email exists: %v", err)
			return ErrEmailExists
		}
		log.Printf("Error inserting user: %v", err)
		return err
	}

	log.Printf("User %s inserted successfully", user.Username)
	return nil
}

func GetUserByID(id int) (*User, error) {
	stmt, err := GetDB().Prepare("SELECT id, username, name, email, currency, password FROM users WHERE id=?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return nil, err
	}
	defer stmt.Close()

	user := &User{}
	err = stmt.QueryRow(id).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		log.Printf("Error retrieving user by ID: %v", err)
		return nil, ErrUserNotFound
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	stmt, err := GetDB().Prepare("SELECT id, username, name, email, currency, password FROM users WHERE username=?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return nil, err
	}
	defer stmt.Close()

	user := &User{}
	err = stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		log.Printf("Error retrieving user by username: %v", err)
		return nil, ErrUserNotFound
	}
	return user, nil
}

func UpdateUser(user *User) error {
	stmt, err := GetDB().Prepare("UPDATE users SET username=?, name=?, email=?, currency=?, password=? WHERE id=?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Name, user.Email, user.Currency, user.Password, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "username") {
			return ErrUsernameTaken
		}
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "email") {
			return ErrEmailExists
		}
		log.Printf("Error updating user: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		log.Println("No user found to update")
		return ErrUserNotFound
	}
	return nil
}

func DeleteUser(id int) error {
	stmt, err := GetDB().Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		log.Println("No user found to delete")
		return ErrUserNotFound
	}
	return nil
}

func UserExists(username string, email string) (bool, error) {
	stmt, err := GetDB().Prepare("SELECT 1 FROM users WHERE username=? OR email=?")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return false, err
	}
	defer stmt.Close()

	var exists int
	err = stmt.QueryRow(username, email).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("Error checking if user exists: %v", err)
		return false, err
	}
	return exists == 1, nil
}
