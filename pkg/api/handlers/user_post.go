package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"log"
	"net/http"
)

func (h *SocialMediaHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding user data: %v", err)
		http.Error(w, fmt.Sprintf("Invalid input: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.DB.Create(&user).Error; err != nil {
		log.Printf("Database insertion error: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
