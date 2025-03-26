package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	v1alpha1 "go-social-media/pkg/types"
	"log"
	"net/http"
	"strconv"

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

	var getPostReply v1alpha1.PostReturnType

	postKey := fmt.Sprintf("post:%d", id)
	viewsKey := fmt.Sprintf("post:%d:views", id)
	likesKey := fmt.Sprintf("post:%d:likes", id)

	log.Printf("[INFO] Attempting to retrieve post - ID: %d, PostKey: %s, ViewsKey: %s", id, postKey, viewsKey)

	// Attempt to retrieve from cache first
	post, err := h.getPostFromCache(postKey, id)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve post - ID: %d, Error: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	views, err := h.incrementPostViews(viewsKey, post.Views)
	if err != nil {
		log.Printf("[ERROR] Failed to increment views - ViewsKey: %s, Error: %v", viewsKey, err)
	} else {
		log.Printf("[INFO] Views incremented - PostID: %d, NewViewCount: %d", id, views)
		post.Views = views
	}

	likes, err := h.getLikesFromCache(likesKey, id)

	getPostReply.Post = *post
	getPostReply.NumberOfLikes = likes

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Keep-Alive", "timeout=5, max=1000")

	if err := json.NewEncoder(w).Encode(getPostReply); err != nil {
		log.Printf("[ERROR] Failed to encode post - PostID: %d, Error: %v", id, err)
	} else {
		log.Printf("[INFO] Successfully sent post response - PostID: %d, Title: %s", post.ID, post.Title)
	}
}

func (h *SocialMediaHandler) getPostFromCache(postKey string, id int) (*models.Post, error) {
	postResult, redisPostErr := h.RedisReader.Get(postKey).Result()
	if redisPostErr == nil {
		log.Printf("[INFO] Cache hit for post - Key: %s", postKey)
		var post models.Post
		if err := json.Unmarshal([]byte(postResult), &post); err != nil {
			log.Printf("[ERROR] Failed to unmarshal cached post - Key: %s, Error: %v", postKey, err)
			return nil, fmt.Errorf("failed to unmarshal cached post: %w", err)
		}
		return &post, nil
	}

	// Cache miss, log and fetch from database
	log.Printf("[INFO] Cache miss for post - Key: %s, Fetching from database", postKey)
	var post models.Post
	result := h.DBReader.Preload("User").First(&post, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("[WARN] Post not found in database - ID: %d", id)
			return nil, fmt.Errorf("post not found")
		}
		log.Printf("[ERROR] Database query failed - ID: %d, Error: %v", id, result.Error)
		return nil, result.Error
	}

	// Attempt to cache the fetched post
	marshalledPost, err := json.Marshal(post)
	if err != nil {
		log.Printf("[WARN] Failed to marshal post for caching - PostID: %d, Error: %v", id, err)
	} else {
		if err := h.RedisReader.Set(postKey, marshalledPost, CACHE_DURATION_VERY_LONG).Err(); err != nil {
			log.Printf("[WARN] Failed to cache post - Key: %s, Error: %v", postKey, err)
		} else {
			log.Printf("[INFO] Post successfully cached - Key: %s, PostID: %d", postKey, id)
		}
	}

	return &post, nil
}

func (h *SocialMediaHandler) getLikesFromCache(likesKey string, id int) (int64, error) {
	var likesCount int64
	likesCountString, redisLikesErr := h.RedisReader.Get(likesKey).Result()
	if redisLikesErr == nil {
		log.Printf("[INFO] Cache hit for post - Key: %s", likesKey)
		if err := json.Unmarshal([]byte(likesCountString), &likesCount); err != nil {
			log.Printf("[ERROR] Failed to unmarshal cached post - Key: %s, Error: %v", likesKey, err)
			return 0, fmt.Errorf("failed to unmarshal cached post: %w", err)
		}
		return likesCount, nil
	}

	// Cache miss, log and fetch from database
	log.Printf("[INFO] Cache miss for likes - Key: %s, Fetching from database", likesKey)
	result := h.DBReader.Model(&models.Like{}).Where("post_id = ?", id).Count(&likesCount)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("[WARN] Post not found in database - ID: %d", id)
			return 0, fmt.Errorf("post not found")
		}
		log.Printf("[ERROR] Database query failed - ID: %d, Error: %v", id, result.Error)
		return 0, result.Error
	}

	// Attempt to cache the fetched post
	if err := h.RedisReader.Set(likesKey, likesCount, CACHE_DURATION_VERY_LONG).Err(); err != nil {
		log.Printf("[WARN] Failed to update views key - Key: %s, Error: %v", likesKey, err)
	}

	return likesCount, nil
}

func (h *SocialMediaHandler) incrementPostViews(viewsKey string, dbPostViews int) (int, error) {
	// Atomically increment views
	viewsCount, err := h.RedisReader.Incr(viewsKey).Result()
	if err != nil {
		log.Printf("[ERROR] Failed to increment views - Key: %s, Error: %v", viewsKey, err)
		return dbPostViews, err
	}

	finalViewCount := max(int(viewsCount), dbPostViews+1)

	if err := h.RedisReader.Set(viewsKey, finalViewCount, CACHE_DURATION_VERY_LONG).Err(); err != nil {
		log.Printf("[WARN] Failed to update views key - Key: %s, Error: %v", viewsKey, err)
	}

	return finalViewCount, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
