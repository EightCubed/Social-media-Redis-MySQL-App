package handlers

import (
	"go-social-media/pkg/database"
	"go-social-media/pkg/models"
	"log"
	"strconv"
	"time"

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

func SyncViewsToDB(db *database.DBConnection, redisClient *redis.Client) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		keys, err := redisClient.Keys("post:*:views").Result()
		if err != nil {
			log.Printf("[ERROR] Failed to fetch Redis keys: %v", err)
			continue
		}

		for _, key := range keys {
			postID, _ := strconv.Atoi(key[5 : len(key)-6])
			views, err := redisClient.Get(key).Int()
			if err != nil {
				log.Printf("[ERROR] Failed to get views for %s: %v", key, err)
				continue
			}

			if views > 0 {
				result := db.GormDBWriter.Model(&models.Post{}).Where("id = ?", postID).Update("views", gorm.Expr("views + ?", views))
				if result.Error != nil {
					log.Printf("[ERROR] Failed to update views in DB: %v", result.Error)
				} else {
					log.Printf("[INFO] Flushed %d views to MySQL for post %d", views, postID)
					redisClient.Del(key)
				}
			}
		}
	}
}
