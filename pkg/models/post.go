package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index"`
	Title   string `gorm:"size:255;not null"`
	Content string `gorm:"type:text;not null"`
	Views   int    `gorm:"default:0"`

	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Comments []Comment `gorm:"foreignKey:PostID"`
	Likes    []Like    `gorm:"foreignKey:PostID"`
}
