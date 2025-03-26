package database

import (
	"fmt"
	"go-social-media/pkg/config"
	"go-social-media/pkg/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DatabaseWriterInit(config config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser,
		config.DBPassword,
		config.DBWriteHost,
		config.DBName)

	log.Printf("Connecting to write database...")

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to open write database connection: %v", err)
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Printf("Failed to get SQL DB from GORM: %v", err)
		return nil, fmt.Errorf("failed to get SQL DB from GORM: %v", err)
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	log.Printf("Write Database connection successful.")

	log.Printf("Creating tables if they do not exist...")
	err = AutoMigrateTables(gormDB)
	if err != nil {
		log.Printf("Failed to create tables: %v", err)
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}
	log.Printf("Tables initialized successfully.")

	return gormDB, nil
}

func DatabaseReaderInit(config config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser,
		config.DBPassword,
		config.DBReadHost,
		config.DBName)

	log.Printf("Connecting to read database...")

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to open read database connection: %v", err)
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Printf("Failed to get SQL DB from GORM: %v", err)
		return nil, fmt.Errorf("failed to get SQL DB from GORM: %v", err)
	}

	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	log.Printf("Read Database connection successful.")

	return gormDB, nil
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
