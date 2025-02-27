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
			"user_id" VARCHAR(36) NOT NULL UNIQUE,
			FOREIGN KEY(user_id) REFERENCES users(uuid)
		);`

	commentsTable := `CREATE TABLE IF NOT EXISTS comments (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"uuid" VARCHAR(36) NOT NULL UNIQUE,
			"content" TEXT NOT NULL,
			"created_at" DATETIME NOT NULL,
			"user_id" VARCHAR(36) NOT NULL UNIQUE,
			"post_id" VARCHAR(36) NOT NULL UNIQUE,
			FOREIGN KEY(user_id) REFERENCES users(uuid),
			FOREIGN KEY(post_id) REFERENCES posts(uuid) ON DELETE CASCADE

		);`

	messagesTable := `CREATE TABLE IF NOT EXISTS messages(
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"uuid" VARCHAR(36) NOT NULL UNIQUE,
			"sender_id" VARCHAR(36) NOT NULL UNIQUE,
			"receiver_id" VARCHAR(36) NOT NULL UNIQUE,
			"content" TEXT NOT NULL,
			"created_at" DATETIME NOT NULL,
			FOREIGN KEY(sender_id) REFERENCES users(uuid),
			FOREIGN KEY(receiver_id) REFERENCES users(uuid)
	)`

	CreateTable(db, usersTable, "users")
	CreateTable(db, postsTable, "posts")
	CreateTable(db, commentsTable, "comments")
	CreateTable(db, messagesTable, "messages")

}

func CreateTable(db *sql.DB, createTableSQL string, tableName string) { //create one table already defined
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Creation table Failed %s : %v", tableName, err)
	}
	fmt.Printf("Table %s allready exist.\n", tableName)
}
