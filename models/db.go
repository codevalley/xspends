package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB
var SQLBuilder squirrel.StatementBuilderType
var logrs = logrus.New()

func GetQueryBuilder() *squirrel.StatementBuilderType {
	return &SQLBuilder
}

func GetDB() *sql.DB {
	return DB
}

func InitDB() {
	var err error
	dsn := os.Getenv("DB_DSN")
	fmt.Println(dsn)
	if dsn == "" {
		log.Fatal("DB_DSN environment variable not set.")
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
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

func GetTransaction(tx ...*sql.Tx) (isExternalTx bool, newTx *sql.Tx, err error) {
	db := GetDB()

	if len(tx) > 0 && tx[0] != nil {
		isExternalTx = true
		newTx = tx[0]
	} else {
		newTx, err = db.Begin()
	}

	return
}
