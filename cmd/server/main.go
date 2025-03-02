package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go-social-media/pkg/api/handlers"
	config "go-social-media/pkg/config"
	database "go-social-media/pkg/database"
)

type App struct {
	DB     *database.DBConnection
	Router *mux.Router
	Config config.Config
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
	socialMediaHandler := handlers.ReturnHandler(a.DB.GormDB)

	apiRouter := a.Router.PathPrefix("/apis/v1").Subrouter()

	// Health endpoints
	apiRouter.HandleFunc("/health", socialMediaHandler.HealthCheck).Methods("GET")

	// User endpoints
	apiRouter.HandleFunc("/user/{id:[0-9]+}", socialMediaHandler.GetUser).Methods("GET")
	apiRouter.HandleFunc("/user", socialMediaHandler.PostUser).Methods("POST")
	// apiRouter.HandleFunc("/user", socialMediaHandler.UpdateUser).Methods("PUT")
	// apiRouter.HandleFunc("/user", socialMediaHandler.PatchUser).Methods("PATCH")
	// apiRouter.HandleFunc("/user", socialMediaHandler.DeleteUser).Methods("DELETE")

	// Post endpoints
	// apiRouter.HandleFunc("/user", socialMediaHandler.DeleteUser).Methods("DELETE")

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
		ServerPort: config.GetEnv("SERVER_PORT", "8080"),
	}

	app := App{Config: config}

	err := app.Initialize()
	if err != nil {
		log.Fatalf("Application failed to initialize: %v", err)
	}
	defer database.DBClose(app.DB)

	app.Run()
}
