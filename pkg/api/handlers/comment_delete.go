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

func (h *SocialMediaHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] DeleteComment handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["comment_id"])
	if err != nil {
		log.Printf("[ERROR] Invalid comment ID: %v", err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	result := h.DBWriter.Delete(&models.Comment{}, commentID)

	if result.Error != nil {
		log.Printf("[ERROR] Database query error: %v", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("[WARN] No comment found with ID: %d", commentID)
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	commentsListKey := fmt.Sprintf("commentlist:%d", commentID)
	res, err := h.RedisReader.Del(commentsListKey).Result()
	if err != nil {
		log.Printf("Failed to delete key: %v", err)
		return
	}
	log.Printf("[INFO] Deleted %d keys", res)

	log.Printf("[INFO] Successfully deleted post with ID: %d", commentID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Comment deleted successfully",
		"post_id": commentID,
	})
}
