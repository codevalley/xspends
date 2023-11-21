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
	SQLBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	return nil
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
