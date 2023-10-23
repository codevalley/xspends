package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func InitializeStore() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(tidb-cluster-tidb.tidb-cluster.svc.cluster.local:4000)/xpends")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	log.Println("Successfully connected to the database")

	// Configure database connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60)
}
