package handlers

import (
	"encoding/json"
	"go-social-media/pkg/models"
	"log"
	"net/http"
)

func (h *SocialMediaHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] ListUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	var users []models.User

	result := h.DBReader.Find(&users)

	if result.Error != nil {
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[INFO] Successfully retrieved list")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
