package models

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
	"xspends/util"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Currency  string    `json:"currency"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

	sid, err := util.GenerateSnowflakeID()
	if err != nil {
		log.Printf("[ERROR] Generating snowflake ID for user: %v", err)
		return err
	}
	user.ID = sid
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	stmt, err := GetDB().Prepare("INSERT INTO users (id, username, name, email, currency, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[ERROR] Preparing insert statement for user: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Username, user.Name, user.Email, user.Currency, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "username") {
			log.Printf("[ERROR] Username taken: %v", err)
			return ErrUsernameTaken
		}
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "email") {
			log.Printf("[ERROR] Email exists: %v", err)
			return ErrEmailExists
		}
		log.Printf("[ERROR] Inserting user: %v", err)
		return err
	}

	log.Printf("User %s inserted successfully", user.Username)
	return nil
}

func GetUserByID(id int64) (*User, error) {
	stmt, err := GetDB().Prepare("SELECT id, username, name, email, currency, password, created_at, updated_at FROM users WHERE id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing select statement for user by ID: %v", err)
		return nil, err
	}
	defer stmt.Close()

	user := &User{}
	err = stmt.QueryRow(id).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		log.Printf("[ERROR] Retrieving user by ID: %v", err)
		return nil, err
	}

	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	stmt, err := GetDB().Prepare("SELECT id, username, name, email, currency, password, created_at, updated_at FROM users WHERE username=?")
	if err != nil {
		log.Printf("[ERROR] Preparing select statement for user by username: %v", err)
		return nil, err
	}
	defer stmt.Close()

	user := &User{}
	err = stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		log.Printf("[ERROR] Retrieving user by username: %v", err)
		return nil, err
	}

	return user, nil
}

func UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()

	stmt, err := GetDB().Prepare("UPDATE users SET username=?, name=?, email=?, currency=?, password=?, updated_at=? WHERE id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing update statement for user: %v", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Name, user.Email, user.Currency, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "username") {
			return ErrUsernameTaken
		}
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "email") {
			return ErrEmailExists
		}
		log.Printf("[ERROR] Updating user: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[ERROR] Fetching rows affected after updating user: %v", err)
		return err
	}
	if rowsAffected == 0 {
		log.Println("[WARNING] No user found to update")
		return ErrUserNotFound
	}
	return nil
}

func DeleteUser(id int64) error {
	stmt, err := GetDB().Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		log.Printf("[ERROR] Preparing delete statement for user: %v", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		log.Printf("[ERROR] Deleting user: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[ERROR] Fetching rows affected after deleting user: %v", err)
		return err
	}
	if rowsAffected == 0 {
		log.Println("[WARNING] No user found to delete")
		return ErrUserNotFound
	}
	return nil
}

func UserExists(username string, email string) (bool, error) {
	stmt, err := GetDB().Prepare("SELECT 1 FROM users WHERE username=? OR email=?")
	if err != nil {
		log.Printf("[ERROR] Preparing select statement to check if user exists: %v", err)
		return false, err
	}
	defer stmt.Close()

	var exists int
	err = stmt.QueryRow(username, email).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("[ERROR] Checking if user exists: %v", err)
		return false, err
	}
	return exists == 1, nil
}
