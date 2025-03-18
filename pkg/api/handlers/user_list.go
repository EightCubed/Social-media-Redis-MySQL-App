package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

func (h *SocialMediaHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] ListUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	res, redisErr := h.RedisReader.Get("userlist").Result()
	if redisErr == redis.Nil {
		log.Printf("[INFO] Cache miss")
	} else if redisErr != nil {
		log.Printf("[ERROR] Failed to get cache: %v", redisErr)
	} else {
		fmt.Printf("[INFO] Cache hit:", res)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(res))
		return
	}

	var users []models.User

	result := h.DBReader.Find(&users)

	if result.Error != nil {
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.RedisReader.Set("userlist", users, 1*time.Minute)

	log.Printf("[INFO] Successfully retrieved list")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
