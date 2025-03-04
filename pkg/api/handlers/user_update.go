package handlers

import (
	"encoding/json"
	"go-social-media/pkg/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *SocialMediaHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] UpdateUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("[ERROR] Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	updates := map[string]interface{}{}
	if userName := vars["username"]; userName != "" {
		updates["username"] = userName
	}
	if email := vars["email"]; email != "" {
		updates["email"] = email
	}
	if password := vars["password"]; password != "" {
		updates["password"] = password
	}

	result := h.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("[WARN] No user found with ID: %d", id)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.Printf("[INFO] Successfully updated user with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User updated successfully",
		"user_id": id,
	})
}
