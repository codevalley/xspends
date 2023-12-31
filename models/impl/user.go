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

*/
package impl

import (
	"context"
	"database/sql"
	"log"

	"strings"
	"time"
	"xspends/models/interfaces"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

// User struct mirrors the users table and satisfies the AuthBoss User interface.

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email already exists")
	ErrUsernameTaken = errors.New("username already exists")
)

type UserModel struct {
	// Any fields if required
}

func (um *UserModel) InsertUser(ctx context.Context, user *interfaces.User, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	if user.Username == "" {
		return errors.New("mandatory field missing: Username")
	}
	if user.Email == "" {
		return errors.New("mandatory field missing: Email")
	}
	if user.Password == "" {
		return errors.New("mandatory field missing: Password")
	}

	var err error
	user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()
	user.ID, err = util.GenerateSnowflakeID()
	if err != nil {
		return errors.Wrap(err, "generating Snowflake ID failed")
	}

	// Build and execute the SQL query using Squirrel
	sqlquery, args, err := squirrel.Insert("users").
		Columns("id", "username", "name", "email", "currency", "password", "created_at", "updated_at").
		Values(user.ID, user.Username, user.Name, user.Email, user.Currency, user.Password, user.CreatedAt, user.UpdatedAt).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "building SQL query for InsertUser failed")
	}

	log.Printf("SQL query: %s\n", sqlquery)
	// Execute the query
	_, err = executor.ExecContext(ctx, sqlquery, args...)
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

	// Handle transaction commitment if it's not external
	commitOrRollback(executor, isExternalTx, err)

	return nil
}

func (um *UserModel) UpdateUser(ctx context.Context, user *interfaces.User, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)
	user.UpdatedAt = time.Now()

	sqlquery, args, err := squirrel.Update("users").
		SetMap(map[string]interface{}{
			"username":   user.Username,
			"name":       user.Name,
			"email":      user.Email,
			"currency":   user.Currency,
			"password":   user.Password,
			"updated_at": user.UpdatedAt,
		}).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "building SQL query for UpdateUser failed")
	}

	_, err = executor.ExecContext(ctx, sqlquery, args...)
	if err != nil {
		return errors.Wrap(err, "updating user failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (um *UserModel) DeleteUser(ctx context.Context, id int64, otx ...*sql.Tx) error {
	isExternalTx, executor := getExecutor(otx...)

	sqlquery, args, err := squirrel.Delete("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return errors.Wrap(err, "building SQL query for DeleteUser failed")
	}

	_, err = executor.ExecContext(ctx, sqlquery, args...)
	if err != nil {
		return errors.Wrap(err, "deleting user failed")
	}

	commitOrRollback(executor, isExternalTx, err)
	return nil
}

func (um *UserModel) GetUserByID(ctx context.Context, id int64, otx ...*sql.Tx) (*interfaces.User, error) {
	_, executor := getExecutor(otx...)

	sqlquery, args, err := squirrel.Select("id", "username", "name", "email", "currency", "password").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "building SQL query for GetUserByID failed")
	}

	user := &interfaces.User{}
	err = executor.QueryRowContext(ctx, sqlquery, args...).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "retrieving user by ID failed")
	}
	return user, nil
}

func (um *UserModel) GetUserByUsername(ctx context.Context, username string, otx ...*sql.Tx) (*interfaces.User, error) {
	_, executor := getExecutor(otx...)

	sqlquery, args, err := squirrel.Select("id", "username", "name", "email", "currency", "password").
		From("users").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "building SQL query for GetUserByUsername failed")
	}

	user := &interfaces.User{}

	err = executor.QueryRowContext(ctx, sqlquery, args...).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Currency, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "retrieving user by username failed")
	}
	return user, nil
}
func (um *UserModel) UserExists(ctx context.Context, username, email string, otx ...*sql.Tx) (bool, error) {
	_, executor := getExecutor(otx...)

	sqlquery, args, err := squirrel.Select("1").
		From("users").
		Where(squirrel.Or{squirrel.Eq{"username": username}, squirrel.Eq{"email": email}}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return false, errors.Wrap(err, "building SQL query for UserExists failed")
	}

	var exists bool
	err = executor.QueryRowContext(ctx, sqlquery, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "checking if user exists failed")
	}

	return exists, nil
}

func (um *UserModel) UserIDExists(ctx context.Context, id int64, otx ...*sql.Tx) (bool, error) {
	_, executor := getExecutor(otx...)

	sqlquery, args, err := squirrel.Select("1").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return false, errors.Wrap(err, "building SQL query for UserIDExists failed")
	}

	var exists bool
	err = executor.QueryRowContext(ctx, sqlquery, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "checking if user ID exists failed")
	}

	return exists, nil
}
