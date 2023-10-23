package models

import (
	"errors"
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

const (
	ErrUserNotFound = "user not found"
	ErrEmailExists  = "email already exists"
)

// InsertUser adds a new user to the database.
func InsertUser(user *User) error {
	_, err := GetDB().Exec("INSERT INTO users (id, username, name, email, currency, password) VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Name, user.Email, user.Currency, user.Password)
	return err
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(id string) (*User, error) {
	row := GetDB().QueryRow("SELECT id, username, name, email, currency, password FROM users WHERE id=?", id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		return nil, errors.New(ErrUserNotFound)
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(username string) (*User, error) {
	row := GetDB().QueryRow("SELECT id, username, name, email, currency, password FROM users WHERE username=?", username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		return nil, errors.New(ErrUserNotFound)
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
	_, err := GetDB().Exec("UPDATE users SET username=?, name=?, email=?, currency=?, password=? WHERE id=?",
		user.Username, user.Name, user.Email, user.Currency, user.Password, user.ID)
	return err
}

// DeleteUser deletes a user from the database.
func DeleteUser(id string) error {
	_, err := GetDB().Exec("DELETE FROM users WHERE id=?", id)
	return err
}
