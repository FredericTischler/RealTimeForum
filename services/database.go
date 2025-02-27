package services

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	err = database.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to SQLite database")
	return database, nil
}
