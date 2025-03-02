package handlers

import (
	"net/http"

	"gorm.io/gorm"
)

type SocialMediaHandler struct {
	DB *gorm.DB
}
type DedicatedComputeHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	PostUser(w http.ResponseWriter, r *http.Request)
	// UpdateUser(w http.ResponseWriter, r *http.Request)
	// PatchUser(w http.ResponseWriter, r *http.Request)
	// DeleteUser(w http.ResponseWriter, r *http.Request)
}

func ReturnHandler(db *gorm.DB) *SocialMediaHandler {
	return &SocialMediaHandler{
		DB: db,
	}
}
