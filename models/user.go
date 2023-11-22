/*
MIT License

Copyright (c) 2023 Narayan Babu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
/*
Package models provides functionality for user management in a Go application.

The User struct represents a user record in the database and implements the AuthBoss User interface.

The package includes functions for inserting, retrieving, updating, and deleting user records in the database. It also includes functions for checking if a user exists and retrieving a user by their ID or username.

Example Usage:

  - Insert a new user:
    user := &User{
    Username: "john_doe",
    Name: "John Doe",
    Email: "john@example.com",
    Currency: "USD",
    Password: "password123",
    }
    err := InsertUser(context.Background(), user)
    if err != nil {
    log.Fatal(err)
    }

  - Retrieve a user by ID:
    retrievedUser, err := GetUserByID(context.Background(), user.ID)
    if err != nil {
    log.Fatal(err)
    }

  - Update a user:
    retrievedUser.Name = "Jane Doe"
    err = UpdateUser(context.Background(), retrievedUser)
    if err != nil {
    log.Fatal(err)
    }

  - Delete a user:
    err = DeleteUser(context.Background(), retrievedUser.ID)
    if err != nil {
    log.Fatal(err)
    }

  - Check if a user exists:
    exists, err := UserExists(context.Background(), "john_doe", "john@example.com")
    if err != nil {
    log.Fatal(err)
    }
    if exists {
    fmt.Println("User exists")
    } else {
    fmt.Println("User does not exist")
    }

  - Check if a user ID exists:
    exists, err := UserIDExists(context.Background(), user.ID)
    if err != nil {
    log.Fatal(err)
    }
    if exists {
    fmt.Println("User ID exists")
    } else {
    fmt.Println("User ID does not exist")
    }
*/
package models

import (
	"context"
	"database/sql"

	"strconv"
	"strings"
	"time"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// User struct mirrors the users table and satisfies the AuthBoss User interface.
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

// Authboss methods
func (u *User) PutPID(pid string) {
	u.ID, _ = strconv.ParseInt(pid, 10, 64)
}

func (u User) GetPID() string {
	return strconv.FormatInt(u.ID, 10)
}

func (u *User) PutPassword(password string) {
	u.Password = password
}

func (u User) GetPassword() string {
	return u.Password
}

func InsertUser(ctx context.Context, user *User, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(ctx, otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	if user.Username == "" {
		return errors.New("mandatory field missing: Username")
	}
	if user.Email == "" {
		return errors.New("mandatory field missing: Email")
	}
	if user.Password == "" {
		return errors.New("mandatory field missing: Password")
	}

	user.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		return errors.Wrap(err, "generating Snowflake ID failed")
	}
	user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()

	sqlquery, args, err := squirrel.Insert("users").
		Columns("id", "username", "name", "email", "currency", "password", "created_at", "updated_at").
		Values(user.ID, user.Username, user.Name, user.Email, user.Currency, user.Password, user.CreatedAt, user.UpdatedAt).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "building SQL query for InsertUser failed")
	}

	_, err = tx.ExecContext(ctx, sqlquery, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			if strings.Contains(err.Error(), "username") {
				return ErrUsernameTaken
			}
			if strings.Contains(err.Error(), "email") {
				return ErrEmailExists
			}
		}
		return errors.Wrap(err, "inserting user failed")
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}

func GetUserByID(ctx context.Context, id int64, otx ...*sql.Tx) (*User, error) {
	isExternalTx, tx, err := GetTransaction(ctx, otx...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting transaction")
	}
	sqlquery, args, err := squirrel.Select("id", "username", "name", "email", "currency", "password").
		From("users").
		Where(squirrel.Eq{"id": id}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "building SQL query for GetUserByID failed")
	}

	row := tx.QueryRowContext(ctx, sqlquery, args...)
	user := &User{}
	err = row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "retrieving user by ID failed")
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return nil, errors.Wrap(err, "committing transaction failed")
		}
	}
	return user, nil
}

func GetUserByUsername(ctx context.Context, username string, otx ...*sql.Tx) (*User, error) {
	isExternalTx, tx, err := GetTransaction(ctx, otx...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting transaction")
	}
	sqlquery, args, err := squirrel.Select("id", "username", "name", "email", "currency", "password").
		From("users").
		Where(squirrel.Eq{"username": username}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "building SQL query for GetUserByUsername failed")
	}

	row := tx.QueryRowContext(ctx, sqlquery, args...)
	user := &User{}
	err = row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "retrieving user by username failed")
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return nil, errors.Wrap(err, "committing transaction failed")
		}
	}
	return user, nil
}

func UpdateUser(ctx context.Context, user *User, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(ctx, otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	user.UpdatedAt = time.Now()

	sql, args, err := squirrel.Update("users").
		SetMap(map[string]interface{}{
			"username":   user.Username,
			"name":       user.Name,
			"email":      user.Email,
			"currency":   user.Currency,
			"password":   user.Password,
			"updated_at": user.UpdatedAt,
		}).
		Where(squirrel.Eq{"id": user.ID}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "building SQL query for UpdateUser failed")
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			if strings.Contains(err.Error(), "username") {
				return ErrUsernameTaken
			}
			if strings.Contains(err.Error(), "email") {
				return ErrEmailExists
			}
		}
		return errors.Wrap(err, "updating user failed")
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}

func DeleteUser(ctx context.Context, id int64, otx ...*sql.Tx) error {

	isExternalTx, tx, err := GetTransaction(ctx, otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	sql, args, err := squirrel.Delete("users").
		Where(squirrel.Eq{"id": id}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "building SQL query for DeleteUser failed")
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "deleting user failed")
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction failed")
		}
	}
	return nil
}

func UserExists(ctx context.Context, username, email string, otx ...*sql.Tx) (bool, error) {
	isExternalTx, tx, err := GetTransaction(ctx, otx...)
	if err != nil {
		return false, errors.Wrap(err, "error getting transaction")
	}

	var exists bool

	sqlquery, args, err := squirrel.Select("1").From("users").Where(squirrel.Or{
		squirrel.Eq{"username": username},
		squirrel.Eq{"email": email},
	}).RunWith(tx).PlaceholderFormat(squirrel.Question).ToSql()

	if err != nil {
		return false, errors.Wrap(err, "building SQL query for UserExists failed")
	}

	err = tx.QueryRowContext(ctx, sqlquery, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "checking if user exists failed")
	}
	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return false, errors.Wrap(err, "committing transaction failed")
		}
	}
	return exists, nil
}

func UserIDExists(ctx context.Context, id int64, otx ...*sql.Tx) (bool, error) {

	isExternalTx, tx, err := GetTransaction(ctx, otx...)
	if err != nil {
		return false, errors.Wrap(err, "error getting transaction")
	}
	var exists bool

	sqlquery, args, err := squirrel.Select("1").From("users").Where(squirrel.Eq{"id": id}).
		RunWith(tx).PlaceholderFormat(squirrel.Question).ToSql()

	if err != nil {
		return false, errors.Wrap(err, "building SQL query for UserIDExists failed")
	}

	err = tx.QueryRowContext(ctx, sqlquery, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "checking if user ID exists failed")
	}

	if !isExternalTx {
		err = tx.Commit()
		if err != nil {
			return false, errors.Wrap(err, "committing transaction failed")
		}
	}

	return exists, nil
}
