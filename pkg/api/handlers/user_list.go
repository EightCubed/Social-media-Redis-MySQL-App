package handlers

import (
	"encoding/json"
	"go-social-media/pkg/models"
	"log"
	"net/http"
)

func (h *SocialMediaHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] ListUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	var users []models.User

	result := h.DB.Find(&users)

	if result.Error != nil {
		log.Println("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("[INFO] Successfully retrieved list")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
