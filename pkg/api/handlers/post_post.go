package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"log"
	"net/http"
)

type PostPostBody struct {
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *SocialMediaHandler) PostPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] PostPost handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	var body PostPostBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("[ERROR] Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Initialize views to 0
	post := &models.Post{
		UserID:  uint(body.UserID),
		Title:   body.Title,
		Content: body.Content,
		Views:   0,
	}

	log.Printf("[INFO] Creating new post - Title: %s, PostID: %d", post.Title, post.UserID)
	if err := h.DBWriter.Create(&post).Error; err != nil {
		log.Printf("[ERROR] Database insertion failed: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	cacheKey := fmt.Sprintf("post:%d", post.ID)
	marshalledPost, err := json.Marshal(post)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
	}
	err = h.RedisReader.Set(cacheKey, marshalledPost, CACHE_DURATION_LONG).Err()
	if err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] Post created successfully - ID: %d, Title: %s", post.ID, post.Title)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}
