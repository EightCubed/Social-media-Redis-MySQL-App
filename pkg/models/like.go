package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	Id        uuid.UUID
	PostId    uuid.UUID
	UserId    uuid.UUID
	CreatedAt time.Time
}
