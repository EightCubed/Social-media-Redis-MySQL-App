package handlers

import (
	"encoding/json"
	"go-social-media/pkg/models"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

func (h *SocialMediaHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] ListUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	res, redisErr := h.RedisReader.Get(USER_LIST_CACHE_KEY).Result()
	if redisErr == redis.Nil {
		log.Printf("[INFO] Cache miss")
	} else if redisErr != nil {
		log.Printf("[ERROR] Failed to get cache: %v", redisErr)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Keep-Alive", "timeout=5, max=1000")
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

	marshalledUser, err := json.Marshal(users)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
	}
	err = h.RedisReader.Set(USER_LIST_CACHE_KEY, marshalledUser, CACHE_DURATION_LONG).Err()
	if err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] Successfully retrieved list")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Keep-Alive", "timeout=5, max=1000")
	json.NewEncoder(w).Encode(users)
}
