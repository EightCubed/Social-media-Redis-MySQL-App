package handlers

import (
	"database/sql"
	"net/http"
)

type SocialMediaHandler struct {
	DB *sql.DB
}

type DedicatedComputeHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	PostUser(w http.ResponseWriter, r *http.Request)
}

func ReturnHandler(db *sql.DB) *SocialMediaHandler {
	return &SocialMediaHandler{
		DB: db,
	}
}
