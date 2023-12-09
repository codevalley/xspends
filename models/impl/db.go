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

package impl

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

func GetDB() *sql.DB {
	return DB
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

func GetQueryBuilder() *squirrel.StatementBuilderType {
	return &SQLBuilder
}

func GetContext() *context.Context {
	ctx := context.Background()
	return &ctx
}

// getExecutor returns the appropriate DBExecutor (transaction or standard DB connection)
func getExecutor(dbService *DBService, otx ...*sql.Tx) (bool, DBExecutor) {
	if len(otx) > 0 && otx[0] != nil {
		// Use the provided transaction and mark it as external
		return true, otx[0]
	}
	// Use the global DB service's executor and mark it as internal
	if dbService == nil {
		return false, GetDBService().Executor
	} else {
		return false, dbService.Executor
	}
}
func getExecutorNew(otx ...*sql.Tx) (bool, DBExecutor) {
	dbService := GetModelsService().DBService

	if len(otx) > 0 && otx[0] != nil {
		return true, otx[0] // Using provided transaction
	} else {
		return false, dbService.Executor // Using global DB service's executor
	}
}

func GetTxn(ctx context.Context, tx ...*sql.Tx) (isExternalTx bool, newTx *sql.Tx, err error) {
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
func commitOrRollback(executor DBExecutor, isExternalTx bool, actionErr error) error {
	if !isExternalTx {
		tx, ok := executor.(*sql.Tx)
		if !ok || actionErr != nil {
			// Rollback the transaction if the type assertion fails or if there's an action error
			if tx != nil {
				if rbErr := tx.Rollback(); rbErr != nil {
					return errors.Wrap(rbErr, "error rolling back transaction")
				}
			}

			if !ok {
				return errors.New("expected *sql.Tx as executor for internal transaction")
			}
			return actionErr
		}

		// If everything went well, commit the transaction
		if err := tx.Commit(); err != nil {
			return errors.Wrap(err, "error committing transaction")
		}
	}
	return actionErr
}
