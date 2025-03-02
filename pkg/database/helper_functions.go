package database

import (
	"database/sql"
	"log"
)

func DBClose(db *sql.DB) {
	if db != nil {
		log.Println("Closing database connection.")
		db.Close()
		log.Println("Database connection closed.")
	}
}
