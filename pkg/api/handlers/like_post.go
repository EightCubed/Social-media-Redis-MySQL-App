package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *SocialMediaHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] LikePost handler called - Method: %s, Path: %s", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["post_id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userIDParam := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.DBReader.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var post models.Post
	if err := h.DBReader.First(&post, postID).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var like models.Like
	like.PostID = uint(postID)
	like.UserID = uint(userID)

	if err := h.DBReader.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err == nil {
		http.Error(w, "User already liked this post", http.StatusOK)
		return
	}

	if err := h.DBWriter.Create(&like).Error; err != nil {
		log.Printf("[ERROR] Database insertion failed: %v", err)
		http.Error(w, "Failed to add like", http.StatusInternalServerError)
		return
	}

	likesKey := fmt.Sprintf("post:%d:likes", postID)

	// Attempt to cache the updated likes
	_, err = h.RedisReader.Incr(likesKey).Result()
	if err != nil {
		log.Printf("[ERROR] Failed to increment views - Key: %s, Error: %v", likesKey, err)
	}

	log.Printf("[INFO] Like added successfully - PostID: %d, UserID: %d", like.PostID, like.UserID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Like added successfully"})
}
