package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

func (h *SocialMediaHandler) ListPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] ListPost handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	res, redisErr := h.RedisReader.Get(POST_LIST_CACHE_KEY).Result()
	if redisErr == redis.Nil {
		log.Printf("[INFO] Cache miss")
	} else if redisErr != nil {
		log.Printf("[ERROR] Failed to get cache: %v", redisErr)
	} else {
		fmt.Printf("[INFO] Cache hit:")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(res))
		return
	}

	var posts []models.Post

	result := h.DBReader.Find(&posts)

	if result.Error != nil {
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	marshalledPost, err := json.Marshal(posts)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
	}
	err = h.RedisReader.Set(POST_LIST_CACHE_KEY, marshalledPost, CACHE_DURATION_LONG).Err()
	if err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] Successfully retrieved list")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
