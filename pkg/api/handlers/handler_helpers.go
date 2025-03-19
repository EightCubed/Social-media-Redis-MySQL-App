package handlers

import (
	"go-social-media/pkg/database"
	"net/http"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type SocialMediaHandler struct {
	DBWriter    *gorm.DB
	DBReader    *gorm.DB
	RedisReader *redis.Client
}
type DedicatedComputeHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	PostUser(w http.ResponseWriter, r *http.Request)
	// UpdateUser(w http.ResponseWriter, r *http.Request)
	// PatchUser(w http.ResponseWriter, r *http.Request)
	// DeleteUser(w http.ResponseWriter, r *http.Request)
}

func ReturnHandler(db *database.DBConnection, redisClient *redis.Client) *SocialMediaHandler {
	return &SocialMediaHandler{
		DBWriter:    db.GormDBWriter,
		DBReader:    db.GormDBReader,
		RedisReader: redisClient,
	}
}
