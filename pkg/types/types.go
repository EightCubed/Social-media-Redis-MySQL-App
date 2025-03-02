package v1alpha1

import (
	"database/sql"
	"go-social-media/pkg/config"

	"github.com/gorilla/mux"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
	Config config.Config
}
