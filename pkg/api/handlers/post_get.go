package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (h *SocialMediaHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] GetPost handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("[ERROR] Invalid Post ID: %v", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("post:%d", id)

	res, redisErr := h.RedisReader.Get(cacheKey).Result()
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

	var post models.Post
	result := h.DBReader.First(&post, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("[WARNING] Post not found - ID: %d", id)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	marshalledPost, err := json.Marshal(post)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
	}
	err = h.RedisReader.Set(cacheKey, marshalledPost, 1*time.Minute).Err()
	if err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] Successfully retrieved Post - ID: %d, Title: %s", post.ID, post.Title)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
