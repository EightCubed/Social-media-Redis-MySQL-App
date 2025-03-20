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

func (h *SocialMediaHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] GetUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("[ERROR] Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("user:%d", id)

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

	var user models.User
	result := h.DBReader.First(&user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("[WARNING] User not found - ID: %d", id)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	marshalledUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
	}
	err = h.RedisReader.Set(cacheKey, marshalledUser, 30*time.Second).Err()
	if err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] Successfully retrieved user - ID: %d, Username: %s", user.ID, user.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
