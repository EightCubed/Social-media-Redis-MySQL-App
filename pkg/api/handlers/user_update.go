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

func (h *SocialMediaHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] UpdateUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		log.Printf("[ERROR] Invalid user ID: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser models.User
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&updatedUser)
	if err != nil {
		log.Printf("[INFO] uni12719 updatedUser: %v", updatedUser)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] uni12719 updatedUser: %v", updatedUser)

	updates := map[string]interface{}{}
	if updatedUser.Username != "" {
		updates["username"] = updatedUser.Username
	}

	if len(updates) == 0 {
		http.Error(w, "No valid fields provided for update", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] uni12719 updates: %v", updates)

	result := h.DBWriter.Model(&models.User{}).Where("id = ?", userID).Updates(updates)

	if result.Error != nil {
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("[WARN] No user found with ID: %d", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var resultUser models.User
	_ = h.DBReader.First(&resultUser, userID)

	cacheKey := fmt.Sprintf("user:%d", resultUser.ID)
	marshalledUser, err := json.Marshal(resultUser)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
	}
	err = h.RedisReader.Set(cacheKey, marshalledUser, CACHE_DURATION_LONG).Err()
	if err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] Successfully updated user with ID: %d", userID)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User updated successfully",
		"user_id": userID,
	})
}
