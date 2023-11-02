package models

import (
	"context"
	"database/sql"

	"strings"
	"time"
	"xspends/util"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
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

func InsertUser(ctx context.Context, user *User, otx ...*sql.Tx) error {
	isExternalTx, tx, err := GetTransaction(otx...)
	if err != nil {
		return errors.Wrap(err, "error getting transaction")
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("mandatory fields missing")
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
	isExternalTx, tx, err := GetTransaction(otx...)
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
	isExternalTx, tx, err := GetTransaction(otx...)
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
	isExternalTx, tx, err := GetTransaction(otx...)
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

	isExternalTx, tx, err := GetTransaction(otx...)
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
	isExternalTx, tx, err := GetTransaction(otx...)
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

	isExternalTx, tx, err := GetTransaction(otx...)
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
