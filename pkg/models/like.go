package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	PostID uint `gorm:"not null;index"`
	UserID uint `gorm:"not null;index"`

	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
