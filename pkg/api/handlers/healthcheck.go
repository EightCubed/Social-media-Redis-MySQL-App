package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func (h SocialMediaHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check endpoint hit.")

	err := h.DB.Ping()
	if err != nil {
		log.Printf("Health check failed: Database connection error: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(fmt.Sprintf("Database connection error: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
