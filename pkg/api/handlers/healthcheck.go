package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func (h *SocialMediaHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := h.DBReader.DB()
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

	pong, err := h.RedisReader.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
