package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (h *SocialMediaHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user := models.User{}
	result := h.DB.First(&user, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Database query error: %v", result.Error), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
