package database

import (
	"gorm.io/gorm"
)

type DBConnection struct {
	GormDBWriter *gorm.DB
	GormDBReader *gorm.DB
}
