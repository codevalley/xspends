package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"xspends/util"

	"github.com/Masterminds/squirrel"
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

	var err error
	user.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		logrs.WithError(err).Error("Generating Snowflake ID failed")
		return util.ErrDatabase // or a more specific error like ErrGeneratingID
	}
	user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()

	query := SQLBuilder.Insert("users").
		Columns("id", "username", "name", "email", "currency", "password", "created_at", "updated_at").
		Values(user.ID, user.Username, user.Name, user.Email, user.Currency, user.Password, user.CreatedAt, user.UpdatedAt)

	_, err = query.RunWith(GetDB()).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			if strings.Contains(err.Error(), "username") {
				logrs.WithError(err).Error("Username taken")
				return ErrUsernameTaken
			}
			if strings.Contains(err.Error(), "email") {
				logrs.WithError(err).Error("Email exists")
				return ErrEmailExists
			}
		}
		logrs.WithError(err).Error("Inserting user failed")
		return err
	}

	logrs.Infof("User %s inserted successfully", user.Username)
	return nil
}

func GetUserByID(id int64) (*User, error) {
	user := &User{}
	query := SQLBuilder.Select("id", "username", "name", "email", "currency", "password").From("users").Where(squirrel.Eq{"id": id})

	err := query.RunWith(GetDB()).QueryRow().Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			logrs.WithError(err).Warn("User not found by ID")
			return nil, ErrUserNotFound
		}
		logrs.WithError(err).Error("Retrieving user by ID failed")
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(username string) (*User, error) {
	user := &User{}

	query := SQLBuilder.Select("id", "username", "name", "email", "currency", "password").From("users").Where(squirrel.Eq{"username": username})
	err := query.RunWith(GetDB()).QueryRow().Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			logrs.WithField("username", username).Warn("User not found by username")
			return nil, ErrUserNotFound
		}
		logrs.WithError(err).Error("Retrieving user by username failed")
		return nil, err
	}

	return user, nil
}

func UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()

	query := SQLBuilder.Update("users").
		SetMap(map[string]interface{}{
			"username":   user.Username,
			"name":       user.Name,
			"email":      user.Email,
			"currency":   user.Currency,
			"password":   user.Password,
			"updated_at": user.UpdatedAt,
		}).
		Where(squirrel.Eq{"id": user.ID})

	result, err := query.RunWith(GetDB()).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			if strings.Contains(err.Error(), "username") {
				return ErrUsernameTaken
			}
			if strings.Contains(err.Error(), "email") {
				return ErrEmailExists
			}
		}
		logrs.WithError(err).Error("Updating user failed")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrs.WithError(err).Error("Fetching rows affected after updating user failed")
		return err
	}
	if rowsAffected == 0 {
		logrs.Warn("No user found to update")
		return ErrUserNotFound
	}

	return nil
}

// DeleteUser deletes a user with the specified ID.
func DeleteUser(id int64) error {
	query := SQLBuilder.Delete("users").Where(squirrel.Eq{"id": id})
	result, err := query.RunWith(GetDB()).Exec()

	if err != nil {
		logrs.WithError(err).Error("Deleting user failed")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrs.WithError(err).Error("Fetching rows affected after deleting user failed")
		return err
	}

	if rowsAffected == 0 {
		logrs.Warn("No user found to delete")
		return ErrUserNotFound
	}

	return nil
}

// UserExists checks if a user with the given username or email exists.
func UserExists(username, email string) (bool, error) {
	var exists bool

	query := SQLBuilder.Select("1").From("users").Where(squirrel.Or{
		squirrel.Eq{"username": username},
		squirrel.Eq{"email": email},
	})
	err := query.RunWith(GetDB()).QueryRow().Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logrs.WithError(err).Error("Checking if user exists failed")
		return false, err
	}

	return exists, nil
}

// UserIDExists checks if a user with the given userID exists.
func UserIDExists(userID int64) (bool, error) {
	var exists bool

	query := SQLBuilder.Select("1").From("users").Where(squirrel.Eq{"id": userID})
	err := query.RunWith(GetDB()).QueryRow().Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		logrs.WithError(err).Error("Checking if user ID exists failed")
		return false, err
	}

	return exists, nil
}
