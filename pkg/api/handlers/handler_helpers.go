package handlers

import (
	"go-social-media/pkg/database"
	"net/http"

	"gorm.io/gorm"
)

type SocialMediaHandler struct {
	DBWriter *gorm.DB
	DBReader *gorm.DB
}
type DedicatedComputeHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	PostUser(w http.ResponseWriter, r *http.Request)
	// UpdateUser(w http.ResponseWriter, r *http.Request)
	// PatchUser(w http.ResponseWriter, r *http.Request)
	// DeleteUser(w http.ResponseWriter, r *http.Request)
}

func ReturnHandler(db *database.DBConnection) *SocialMediaHandler {
	return &SocialMediaHandler{
		DBWriter: db.GormDBWriter,
		DBReader: db.GormDBReader,
	}
}
