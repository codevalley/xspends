package models

import (
	"database/sql"
	"errors"
	"strings"
)

// User represents the user of the application.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Currency string `json:"currency"`
	Password string `json:"-"`
}

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email already exists")
	ErrUsernameTaken = errors.New("username already exists") // New error definition
)

// ... (rest of the functions remain unchanged)

// InsertUser adds a new user to the database.
func InsertUser(user *User) error {
	_, err := GetDB().Exec("INSERT INTO users (id, username, name, email, currency, password) VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Name, user.Email, user.Currency, user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "username") { // This is a generic check; adjust based on your database error message
			return ErrUsernameTaken
		}
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "email") {
			return ErrEmailExists
		}
	}
	return err
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(id string) (*User, error) {
	row := GetDB().QueryRow("SELECT id, username, name, email, currency, password FROM users WHERE id=?", id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(username string) (*User, error) {
	row := GetDB().QueryRow("SELECT id, username, name, email, currency, password FROM users WHERE username=?", username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// EmailExists checks if an email already exists in the database.
func EmailExists(email string) bool {
	row := GetDB().QueryRow("SELECT 1 FROM users WHERE email=?", email)

	var exists bool
	err := row.Scan(&exists)
	return err == nil && exists
}

// UpdateUser updates the user details in the database.
func UpdateUser(user *User) error {
	result, err := GetDB().Exec("UPDATE users SET username=?, name=?, email=?, currency=?, password=? WHERE id=?",
		user.Username, user.Name, user.Email, user.Currency, user.Password, user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "username") {
			return ErrUsernameTaken
		}
		if strings.Contains(err.Error(), "duplicate") && strings.Contains(err.Error(), "email") {
			return ErrEmailExists
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func DeleteUser(id string) error {
	result, err := GetDB().Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// UserExists checks if a user with the provided username or email exists.
func UserExists(username string, email string) (bool, error) {
	var exists bool
	row := GetDB().QueryRow("SELECT 1 FROM users WHERE username=? OR email=?", username, email)

	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User does not exist
		}
		return false, err // Some other database error occurred
	}

	return true, nil // User exists
}
