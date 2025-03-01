package database

import (
	"database/sql"
	"fmt"
	"go-social-media/pkg/config"
	"log"
)

func DatabaseInit(connectionString string, config config.Config) (*sql.DB, error) {
	log.Println("Connecting to database...")
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	log.Println("Database connection successful.")

	log.Println("Creating tables if they do not exist...")
	err = CreateTables(db, config)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}
	log.Println("Tables initialized successfully.")

	return db, nil
}
