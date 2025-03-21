package repositories

import "database/sql"

type MessageRepository struct {
	DB *sql.DB
}
