package handlers

import (
	"go-social-media/pkg/database"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type SocialMediaHandler struct {
	DBWriter    *gorm.DB
	DBReader    *gorm.DB
	RedisReader *redis.Client
}

func ReturnHandler(db *database.DBConnection, redisClient *redis.Client) *SocialMediaHandler {
	return &SocialMediaHandler{
		DBWriter:    db.GormDBWriter,
		DBReader:    db.GormDBReader,
		RedisReader: redisClient,
	}
}
