package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"go-social-media/pkg/api/handlers"
	config "go-social-media/pkg/config"
	database "go-social-media/pkg/database"
)

type App struct {
	DB      *sql.DB
	Router  *mux.Router
	Handler handlers.SocialMediaHandler
	Config  config.Config
}

func (a *App) Initialize() error {
	log.Println("Initializing application...")

	var err error
	a.DB, err = database.DatabaseInit(a.Config)
	if err != nil {
		return err
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	log.Println("Application initialization completed successfully.")
	return nil
}

func (a *App) initializeRoutes() {
	socialMediaHandler := handlers.ReturnHandler(a.Handler.DB)

	apiRouter := a.Router.PathPrefix("/apis/v1").Subrouter()

	// Health endpoints
	apiRouter.HandleFunc("/health", socialMediaHandler.HealthCheck).Methods("GET")

	// User endpoints
	apiRouter.HandleFunc("/user", socialMediaHandler.GetUser).Methods("GET")
	apiRouter.HandleFunc("/user", socialMediaHandler.HealthCheck).Methods("POST")
	apiRouter.HandleFunc("/user", socialMediaHandler.HealthCheck).Methods("PUT")
	apiRouter.HandleFunc("/user", socialMediaHandler.HealthCheck).Methods("PATCH")
	apiRouter.HandleFunc("/user", socialMediaHandler.HealthCheck).Methods("DELETE")

	log.Println("API routes initialized.")
}

func (a *App) Run() {
	log.Printf("Starting server on port %s...\n", a.Config.ServerPort)
	err := http.ListenAndServe(":"+a.Config.ServerPort, a.Router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func main() {
	log.Println("Reading environment variables...")

	config := config.Config{
		DBHost:     config.GetEnv("DB_HOST", "mysql.default.svc.cluster.local"),
		DBUser:     config.GetEnv("DB_USER", "root"),
		DBPassword: config.GetEnv("DB_PASSWORD", "rootpassword"),
		DBName:     config.GetEnv("DB_NAME", "social_media_app"),
		ServerPort: config.GetEnv("SERVER_PORT", "3306"),
	}

	app := App{Config: config}

	err := app.Initialize()
	if err != nil {
		log.Fatalf("Application failed to initialize: %v", err)
	}
	defer database.DBClose(app.DB)

	app.Run()
}
