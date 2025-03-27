package models

import (
	"gorm.io/gorm"
)

type Login struct {
	gorm.Model
	PasswordHash string `gorm:"size:255;not null"`
	AccessToken  string `gorm:"type:text"`
	RefreshToken string `gorm:"type:text"`

	UserID *uint `gorm:"uniqueIndex"`
	User   User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
