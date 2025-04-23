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

func (h *SocialMediaHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] UpdatePost handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["post_id"])
	if err != nil {
		log.Printf("[ERROR] Invalid post ID: %v", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var updatedPost models.Post
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&updatedPost)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updates := map[string]interface{}{}
	if updatedPost.Title != "" {
		updates["title"] = updatedPost.Title
	}
	if updatedPost.Content != "" {
		updates["content"] = updatedPost.Content
	}

	if len(updates) == 0 {
		log.Print("[ERROR] Nothing to update")
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	result := h.DBWriter.Model(&models.Post{}).Where("id = ?", postID).Updates(updates)

	if result.Error != nil {
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("[WARN] No post found with ID: %d", postID)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var resultPost models.Post
	result = h.DBReader.First(&resultPost, postID)

	cacheKey := fmt.Sprintf("post:%d", resultPost.ID)
	marshalledPost, err := json.Marshal(resultPost)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
	}
	err = h.RedisReader.Set(cacheKey, marshalledPost, CACHE_DURATION_LONG).Err()
	if err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] Successfully updated post with ID: %d", postID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post updated successfully",
		"post_id": postID,
	})
}
