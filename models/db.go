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

package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	maxIdleConn        = 25
	maxOpenConn        = 25
	maxConnLifetimeMin = 5 * time.Minute
)

var DB *sql.DB
var SQLBuilder squirrel.StatementBuilderType

var dbService *DBService // dbService will hold the instance of DBService
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type DBService struct {
	Executor DBExecutor
}

func (db *DBService) execQuery(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.Executor.ExecContext(ctx, query, args...)
}

func (db *DBService) execQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.Executor.QueryRowContext(ctx, query, args...)
}

func (db *DBService) execQueryRows(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.Executor.QueryContext(ctx, query, args...)
}

func GetQueryBuilder() *squirrel.StatementBuilderType {
	return &SQLBuilder
}

func GetDB() *sql.DB {
	return DB
}

func GetContext() *context.Context {
	ctx := context.Background()
	return &ctx
}

func InitDB() error {
	var err error
	dsn := os.Getenv("DB_DSN")
	fmt.Println(dsn)
	if dsn == "" {
		return errors.New("DB_DSN environment variable not set.")
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return errors.Wrap(err, "Error initializing database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = DB.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, "Error connecting to the database")
	}
	log.Println("Successfully connected to the database")

	// Configure database connection pool
	maxOpenConns := os.Getenv("DB_MAX_OPEN_CONNS")
	if maxOpenConns != "" {
		maxOpen, err := strconv.Atoi(maxOpenConns)
		if err == nil {
			DB.SetMaxOpenConns(maxOpen)
		} else {
			DB.SetMaxOpenConns(maxOpenConn)
		}
	}

	maxIdleConns := os.Getenv("DB_MAX_IDLE_CONNS")
	if maxIdleConns != "" {
		maxIdle, err := strconv.Atoi(maxIdleConns)
		if err == nil {
			DB.SetMaxIdleConns(maxIdle)
		} else {
			DB.SetMaxIdleConns(maxIdleConn)
		}
	}

	connMaxLifetime := os.Getenv("DB_CONN_MAX_LIFETIME")
	if connMaxLifetime != "" {
		connLifetime, err := time.ParseDuration(connMaxLifetime)
		if err == nil {
			DB.SetConnMaxLifetime(connLifetime)
		} else {
			DB.SetConnMaxLifetime(maxConnLifetimeMin)
		}
	}

	dbService = &DBService{
		Executor: DB,
	}

	SQLBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	return nil
}

// GetDBService provides access to the initialized DBService.
func GetDBService() *DBService {
	return dbService
}

// CloseDB safely closes the database connection.
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func GetTransaction(ctx context.Context, tx ...*sql.Tx) (isExternalTx bool, newTx *sql.Tx, err error) {
	db := GetDB()

	if len(tx) > 0 && tx[0] != nil {
		isExternalTx = true
		newTx = tx[0]
	} else {
		newTx, err = db.BeginTx(ctx, nil)
		if err != nil {
			err = errors.Wrap(err, "Error starting a new transaction")
			return
		}
	}

	return
}
