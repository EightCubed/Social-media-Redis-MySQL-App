package database

import (
	"gorm.io/gorm"
)

type DBConnection struct {
	GormDB *gorm.DB
}
