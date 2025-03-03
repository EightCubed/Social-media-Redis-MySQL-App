package database

import (
	"log"
)

func DBClose(conn *DBConnection) {
	if conn != nil && conn.GormDB != nil {
		sqlDB, err := conn.GormDB.DB()
		if err != nil {
			log.Printf("Error getting SQL DB from GORM: %v", err)
			return
		}

		log.Printf("Closing database connection.")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Printf("Database connection closed successfully.")
		}
	}
}
