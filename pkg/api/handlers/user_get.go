package handlers

import (
	"context"
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
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log.Printf("[INFO] GetUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		log.Printf("[ERROR] Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("user:%d", userID)

	res, redisErr := h.RedisReader.WithContext(ctx).Get(cacheKey).Result()
	if redisErr == redis.Nil {
		log.Printf("[INFO] Cache miss")
	} else if redisErr != nil {
		select {
		case <-ctx.Done():
			log.Printf("[ERROR] Redis operation cancelled: %v", ctx.Err())
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
			return
		default:
			log.Printf("[ERROR] Failed to get cache: %v", redisErr)
		}
	} else {
		log.Printf("[INFO] Cache hit")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(res))
		return
	}

	var user models.User
	result := h.DBReader.WithContext(ctx).First(&user, userID)

	if result.Error != nil {
		select {
		case <-ctx.Done():
			log.Printf("[ERROR] Database operation cancelled: %v", ctx.Err())
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
			return
		default:
			if result.Error == gorm.ErrRecordNotFound {
				log.Printf("[WARNING] User not found - ID: %d", userID)
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			log.Printf("[ERROR] Database query error: %v", result.Error)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	marshalledUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
		http.Error(w, "Error processing user data", http.StatusInternalServerError)
		return
	}

	err = h.RedisReader.WithContext(ctx).Set(cacheKey, marshalledUser, CACHE_DURATION_LONG).Err()
	if err != nil {
		select {
		case <-ctx.Done():
			log.Printf("[ERROR] Cache set cancelled: %v", ctx.Err())
		default:
			log.Printf("[ERROR] Cache set error: %v", err)
		}
	}

	log.Printf("[INFO] Successfully retrieved user - ID: %d, Username: %s", user.ID, user.Username)

	select {
	case <-ctx.Done():
		log.Printf("[ERROR] Context done before sending response")
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Keep-Alive", "timeout=5, max=1000")
		json.NewEncoder(w).Encode(user)
	}
}
