package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	PostID uint `gorm:"not null;index;uniqueIndex:idx_post_user,priority:1"`
	UserID uint `gorm:"not null;index;uniqueIndex:idx_post_user,priority:2"`

	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
