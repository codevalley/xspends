package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
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

func InitDB() {
	var err error
	dsn := os.Getenv("DB_DSN")
	fmt.Println(dsn)
	if dsn == "" {
		err = errors.New("DB_DSN environment variable not set.")
		log.Fatal(err)
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		err = errors.Wrap(err, "Error initializing database")
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = DB.PingContext(ctx)
	if err != nil {
		err = errors.Wrap(err, "Error connecting to the database")
		log.Fatal(err)
	}
	log.Println("Successfully connected to the database")

	// Configure database connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute) // Adjusted to use the time package for clarity
	SQLBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
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
