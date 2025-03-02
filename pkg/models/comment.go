package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Id        uuid.UUID
	PostId    uuid.UUID
	UserId    uuid.UUID
	Content   string
	CreatedAt time.Time
}
