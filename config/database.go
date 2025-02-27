package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func BDD() {
	db, err := sql.Open("sqlite3", "furom.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	DefineTables(db)
}

func DefineTables(db *sql.DB) {
	usersTable := `CREATE TABLE IF NOT EXISTS users (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"uuid" VARCHAR(36) NOT NULL UNIQUE,
			"user_name" VARCHAR(255) NOT NULL UNIQUE,
			"email" VARCHAR(255) NOT NULL UNIQUE,
			"password" VARCHAR(72) NOT NULL,
			"first_name" VARCHAR(255) NOT NUL,
			"last_name" VARCHAR(255) NOT NUL,
			"age" INTEGER NOT NULL,
			"gender" VARCHAR(10),
			"created_at" DATETIME NOT NUL
		);`

	CreateTable(db, usersTable, "users")

}

func CreateTable(db *sql.DB, createTableSQL string, tableName string) { //create one table already defined
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Creation table Failed %s : %v", tableName, err)
	}
	fmt.Printf("Table %s allready exist.\n", tableName)
}
