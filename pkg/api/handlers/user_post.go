package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"log"
	"net/http"
	"time"
)

func (h *SocialMediaHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] PostUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("[ERROR] Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] Creating new user - Username: %s, Email: %s", user.Username, user.Email)
	if err := h.DBWriter.Create(&user).Error; err != nil {
		log.Printf("[ERROR] Database insertion failed: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	cacheKey := fmt.Sprintf("user:%d", user.ID)
	marshalledUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
	}
	err = h.RedisReader.Set(cacheKey, marshalledUser, 30*time.Second).Err()
	if err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] User created successfully - ID: %d, Username: %s", user.ID, user.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
