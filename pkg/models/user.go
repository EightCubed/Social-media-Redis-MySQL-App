package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:50;unique;not null"`
	Email    string `gorm:"size:100;unique;not null"`

	LoginID *uint  `gorm:"uniqueIndex"`
	Login   *Login `gorm:"foreignKey:LoginID"`
	Posts   []Post `gorm:"foreignKey:UserID"`
}
