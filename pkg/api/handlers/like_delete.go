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

func (h *SocialMediaHandler) LikeDelete(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] LikeDelete handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["post_id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userIDParam := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil || userID == 0 {
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

	like := &models.Like{
		PostID: uint(postID),
		UserID: uint(userID),
	}

	if err := h.DBReader.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err != nil {
		log.Printf("[WARN] Like not found for PostID: %d by UserID: %d", postID, userID)
		http.Error(w, "Like not found", http.StatusNotFound)
		return
	}

	result := h.DBWriter.Unscoped().Delete(&models.Like{}, like)

	if result.Error != nil {
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("[WARN] No post found with ID: %d", like.PostID)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	likesKey := fmt.Sprintf("post:%d:likes", postID)

	// Attempt to cache the updated likes
	_, err = h.RedisReader.Decr(likesKey).Result()
	if err != nil {
		log.Printf("[ERROR] Failed to increment views - Key: %s, Error: %v", likesKey, err)
	}

	log.Printf("[INFO] Successfully deleted like from PostID: %d", like.PostID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Like deleted successfully",
		"post_id": like.PostID,
	})
}
