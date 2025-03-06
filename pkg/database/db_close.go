package database

import (
	"log"
)

func DBClose(conn *DBConnection) {
	if conn != nil && conn.GormDBWriter != nil {
		sqlDB, err := conn.GormDBWriter.DB()
		if err != nil {
			log.Printf("Error getting SQL DB from GORM: %v", err)
			return
		}

		log.Printf("Closing Writer database connection.")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing Writer database connection: %v", err)
		} else {
			log.Printf("Writer Database connection closed successfully.")
		}
	}

	if conn != nil && conn.GormDBReader != nil {
		sqlDB, err := conn.GormDBReader.DB()
		if err != nil {
			log.Printf("Error getting SQL DB from GORM: %v", err)
			return
		}

		log.Printf("Closing Reader database connection.")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing Reader database connection: %v", err)
		} else {
			log.Printf("Reader Database connection closed successfully.")
		}
	}
}
