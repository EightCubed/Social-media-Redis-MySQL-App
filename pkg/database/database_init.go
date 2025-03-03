package database

import (
	"fmt"
	"go-social-media/pkg/config"
	"go-social-media/pkg/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DatabaseInit(config config.Config) (*DBConnection, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBName)

	log.Printf("Connecting to database...")

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Printf("Failed to get SQL DB from GORM: %v", err)
		return nil, fmt.Errorf("failed to get SQL DB from GORM: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	log.Printf("Database connection successful.")

	log.Printf("Creating tables if they do not exist...")
	err = AutoMigrateTables(gormDB)
	if err != nil {
		log.Printf("Failed to create tables: %v", err)
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}
	log.Printf("Tables initialized successfully.")

	return &DBConnection{
		GormDB: gormDB,
	}, nil
}

func AutoMigrateTables(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto-migrate tables: %v", err)
	}

	log.Printf("✅ Tables migrated successfully!")
	return nil
}
