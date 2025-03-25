package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"log"
	"net/http"
	"strconv"

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

	var post models.Post

	postKey := fmt.Sprintf("post:%d", id)
	viewsKey := fmt.Sprintf("post:%d:views", id)

	postResult, redisPostErr := h.RedisReader.Get(postKey).Result()
	if redisPostErr == redis.Nil || redisPostErr != nil {
		if redisPostErr == redis.Nil {
			log.Printf("[INFO] Cache miss")
		} else {
			log.Printf("[ERROR] Failed to get cache: %v", redisPostErr)
		}

		var result *gorm.DB
		result = h.DBReader.First(&post, id)

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

		var viewCounter int
		viewsResult, redisViewErr := h.RedisReader.Get(viewsKey).Result()
		if redisViewErr == redis.Nil {
			log.Printf("[INFO] Cache miss")
		} else if redisViewErr != nil {
			log.Printf("[ERROR] Failed to get cache: %v", redisViewErr)
		} else {
			if err := json.Unmarshal([]byte(viewsResult), &viewCounter); err != nil {
				log.Printf("[ERROR] Failed to unmarshal cached post: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			} else {
				post.Views = viewCounter
			}
		}

		// Incrementing post views
		post.Views += 1

		// Caching post
		marshalledPost, err := json.Marshal(post)
		if err != nil {
			log.Printf("[ERROR] Marshal error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		err = h.RedisReader.Set(postKey, marshalledPost, CACHE_DURATION_LONG).Err()
		if err != nil {
			log.Printf("[ERROR] Cache set error: %v", err)
		}

		// Caching updated views
		err = h.RedisReader.Set(viewsKey, strconv.Itoa(post.Views), CACHE_DURATION_LONG).Err()
		if err != nil {
			log.Printf("[ERROR] Failed to set views in Redis: %v", err)
		}
	} else {
		log.Printf("[INFO] Cache hit")
		if err := json.Unmarshal([]byte(postResult), &post); err != nil {
			log.Printf("[ERROR] Failed to unmarshal cached post: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var viewCounter int
		viewsResult, redisViewErr := h.RedisReader.Get(viewsKey).Result()
		if redisViewErr == redis.Nil {
			log.Printf("[INFO] Cache miss")
		} else if redisViewErr != nil {
			log.Printf("[ERROR] Failed to get cache: %v", redisViewErr)
		} else {
			if err := json.Unmarshal([]byte(viewsResult), &viewCounter); err != nil {
				log.Printf("[ERROR] Failed to unmarshal cached post: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			} else {
				post.Views = viewCounter
			}
		}

		// Incrementing post views
		post.Views += 1

		// Caching updated views
		err = h.RedisReader.Set(viewsKey, strconv.Itoa(post.Views), CACHE_DURATION_LONG).Err()
		if err != nil {
			log.Printf("[ERROR] Failed to set views in Redis: %v", err)
		}
	}

	log.Printf("[INFO] Successfully retrieved Post - ID: %d, Title: %s, Views: %d", post.ID, post.Title, post.Views)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
