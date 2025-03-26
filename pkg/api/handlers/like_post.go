package handlers

import (
	"encoding/json"
	"go-social-media/pkg/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *SocialMediaHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] LikePost handler called - Method: %s, Path: %s", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
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

	var like models.Like
	like.PostID = uint(postID)
	like.UserID = uint(userID)
	if err := h.DBWriter.Create(&like).Error; err != nil {
		log.Printf("[ERROR] Database insertion failed: %v", err)
		http.Error(w, "Failed to add like", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Like added successfully - PostID: %d, UserID: %s", like.ID, like.UserID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}
