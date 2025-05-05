package database

import (
	"fmt"
	"go-social-media/pkg/config"
	"go-social-media/pkg/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	sqlDB.SetMaxOpenConns(300)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

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

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Printf("Failed to open read database connection: %v", err)
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Printf("Failed to get SQL DB from GORM: %v", err)
		return nil, fmt.Errorf("failed to get SQL DB from GORM: %v", err)
	}

	sqlDB.SetMaxOpenConns(300)
	sqlDB.SetMaxIdleConns(100)
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
	log.Println("Migrating User table first...")
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate User table: %v", err)
	}

	log.Println("Migrating Login table...")
	if err := db.AutoMigrate(&models.Login{}); err != nil {
		return fmt.Errorf("failed to migrate Login table: %v", err)
	}

	log.Println("Migrating remaining tables...")
	if err := db.AutoMigrate(
		&models.Post{},
		&models.Comment{},
		&models.Like{},
	); err != nil {
		return fmt.Errorf("failed to migrate remaining tables: %v", err)
	}

	log.Printf("✅ Tables migrated successfully!")
	return nil
}
