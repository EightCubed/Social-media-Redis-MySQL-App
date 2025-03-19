package handlers

import (
	"encoding/json"
	"go-social-media/pkg/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *SocialMediaHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] DeletePost handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("[ERROR] Invalid post ID: %v", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	result := h.DBWriter.Delete(&models.Post{}, id)

	if result.Error != nil {
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("[WARN] No post found with ID: %d", id)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	log.Printf("[INFO] Successfully deleted post with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post deleted successfully",
		"post_id": id,
	})
}
