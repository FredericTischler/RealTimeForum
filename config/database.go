package config

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

func CreateDataBase() {
	db, err := sql.Open("sqlite3", "forum.db")
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
			"first_name" VARCHAR(255) NOT NULL,
			"last_name" VARCHAR(255) NOT NULL,
			"age" INTEGER NOT NULL,
			"gender" VARCHAR(10),
			"created_at" DATETIME NOT NULL
		);`

	postsTable := `CREATE TABLE IF NOT EXISTS posts (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"uuid" VARCHAR(36) NOT NULL UNIQUE,
			"title" VARCHAR(255) NOT NULL,
			"content" TEXT NOT NULL,
			"category" VARCHAR(30) NOT NULL,
			"created_at" DATETIME NOT NULL,
			"update_at" DATETIME,
			"user_id" VARCHAR(36) NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(uuid) ON DELETE CASCADE
		);`

	commentsTable := `CREATE TABLE IF NOT EXISTS comments (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"uuid" VARCHAR(36) NOT NULL UNIQUE,
			"content" TEXT NOT NULL,
			"created_at" DATETIME NOT NULL,
			"user_id" VARCHAR(36) NOT NULL,
			"post_id" VARCHAR(36) NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(uuid) ON DELETE CASCADE,
			FOREIGN KEY(post_id) REFERENCES posts(uuid) ON DELETE CASCADE

		);`

	messagesTable := `CREATE TABLE IF NOT EXISTS messages(
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"uuid" VARCHAR(36) NOT NULL,
			"sender_id" VARCHAR(36) NOT NULL,
			"receiver_id" VARCHAR(36) NOT NULL,
			"content" TEXT NOT NULL,
			"created_at" DATETIME NOT NULL,
			FOREIGN KEY(sender_id) REFERENCES users(uuid) ON DELETE CASCADE,
			FOREIGN KEY(receiver_id) REFERENCES users(uuid) ON DELETE CASCADE
	)`

	sessionsTable := `CREATE TABLE IF NOT EXISTS sessions(
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"uuid" VARCHAR(36) NOT NULL UNIQUE,
		"user_id" VARCHAR(36) NOT NULL UNIQUE,
		"token" VARCHAR(36) NOT NULL UNIQUE,
		"created_at" DATETIME NOT NULL,
		"expires_at" DATETIME,
		FOREIGN KEY(user_id) REFERENCES users(uuid) ON DELETE CASCADE
	)`

	CreateTable(db, usersTable, "users")
	CreateTable(db, postsTable, "posts")
	CreateTable(db, commentsTable, "comments")
	CreateTable(db, messagesTable, "messages")
	CreateTable(db, sessionsTable, "sessions")
}

func CreateTable(db *sql.DB, createTableSQL string, tableName string) {
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Creation table Failed %s : %v", tableName, err)
	}
	fmt.Printf("Table %s already exist.\n", tableName)
}
