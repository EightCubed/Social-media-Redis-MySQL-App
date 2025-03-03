package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func (h *SocialMediaHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := h.DB.DB()
	if err != nil {
		log.Printf("Health check failed: Unable to get SQL DB: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(fmt.Sprintf("Database connection error: %v", err)))
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("Health check failed: Database connection error: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(fmt.Sprintf("Database connection error: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
