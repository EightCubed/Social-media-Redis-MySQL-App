package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:50;unique;not null"`
	Email    string `gorm:"size:100;unique;not null"`
	Password string `gorm:"size:255;not null"`
	// Posts field is just a GORM association, not an actual DB column
	Posts []Post `gorm:"foreignKey:UserID"`
}
