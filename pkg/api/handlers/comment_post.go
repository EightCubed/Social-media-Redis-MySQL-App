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

type CommentBody struct {
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}

func (h *SocialMediaHandler) PostComment(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] PostComment handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["post_id"])
	if err != nil {
		log.Printf("[ERROR] Invalid Post ID: %v", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var commentBody CommentBody
	if err := json.NewDecoder(r.Body).Decode(&commentBody); err != nil {
		log.Printf("[ERROR] Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	comment := &models.Comment{
		PostID:  uint(postID),
		UserID:  uint(commentBody.UserID),
		Content: commentBody.Content,
	}

	log.Printf("[INFO] Adding new comment - PostID: %d, UserID: %d", comment.PostID, comment.UserID)
	if err := h.DBWriter.Create(&comment).Error; err != nil {
		log.Printf("[ERROR] Database insertion failed: %v", err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	commentsListKey := fmt.Sprintf("commentlist:%d", postID)
	result, err := h.RedisReader.Del(commentsListKey).Result()
	if err != nil {
		log.Printf("Failed to delete key: %v", err)
		return
	}
	log.Printf("[INFO] Deleted %d keys", result)

	log.Printf("[INFO] Comment added successfully - PostID: %d, UserID: %s", comment.PostID, comment.UserID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"})
}
