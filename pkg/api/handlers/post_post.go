package handlers

import (
	"encoding/json"
	"go-social-media/pkg/models"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type PostPostBody struct {
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *SocialMediaHandler) PostPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] PostPost handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	var body PostPostBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		log.Printf("[ERROR] Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if body.Title == "" || body.Content == "" {
		log.Println("[ERROR] Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
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
		var user models.User
		err = h.DBReader.First(&user, post.UserID).Error
		if err == gorm.ErrRecordNotFound {
			log.Printf("[WARNING] User not found - ID: %d", post.UserID)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			log.Printf("[ERROR] Database insertion failed: %v", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("[INFO] Post created successfully - ID: %d, Title: %s", post.ID, post.Title)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}
