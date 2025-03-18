package main

import (
	"go-social-media/pkg/api/handlers"
	config "go-social-media/pkg/config"
	database "go-social-media/pkg/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	DB     *database.DBConnection
	Router *mux.Router
	Config config.Config
}

func (a *App) Initialize() error {
	log.Printf("Initializing application...")

	a.DB = &database.DBConnection{}

	var err error
	a.DB.GormDBReader, err = database.DatabaseReaderInit(a.Config)
	if err != nil {
		return err
	}
	a.DB.GormDBWriter, err = database.DatabaseWriterInit(a.Config)
	if err != nil {
		return err
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	log.Printf("Application initialization completed successfully.")
	return nil
}

func (a *App) initializeRoutes() {
	socialMediaHandler := handlers.ReturnHandler(a.DB)

	apiRouter := a.Router.PathPrefix("/apis/v1").Subrouter()

	// Health endpoints
	apiRouter.HandleFunc("/health", socialMediaHandler.HealthCheck).Methods("GET")

	// User endpoints
	apiRouter.HandleFunc("/user/{id:[0-9]+}", socialMediaHandler.GetUser).Methods("GET")
	apiRouter.HandleFunc("/user", socialMediaHandler.PostUser).Methods("POST")
	apiRouter.HandleFunc("/user", socialMediaHandler.ListUser).Methods("GET")
	apiRouter.HandleFunc("/user/{id:[0-9]+}", socialMediaHandler.UpdateUser).Methods("PATCH")
	apiRouter.HandleFunc("/user/{id:[0-9]+}", socialMediaHandler.DeleteUser).Methods("DELETE")

	// Metrics endpoint
	// a.Router.Handle("/metrics", promhttp.Handler())

	log.Printf("API routes initialized.")
}

func (a *App) Run() {
	log.Printf("Starting server on port %s...\n", a.Config.ServerPort)
	err := http.ListenAndServe(":"+a.Config.ServerPort, a.Router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func main() {
	// logger, _ := zap.NewProduction()
	// defer logger.Sync()

	log.Printf("Reading environment variables...")

	config := config.Config{
		DBWriteHost: config.GetEnv("DB_WRITE_HOST", "mysql-primary.default.svc.cluster.local"),
		DBReadHost:  config.GetEnv("DB_READ_HOST", "mysql-replica.default.svc.cluster.local"),
		DBUser:      config.GetEnv("DB_USER", "root"),
		DBPassword:  config.GetEnv("DB_PASSWORD", "rootpassword"),
		DBName:      config.GetEnv("DB_NAME", "social_media_app"),
		ServerPort:  config.GetEnv("SERVER_PORT", "8080"),
	}

	app := App{Config: config}

	err := app.Initialize()
	if err != nil {
		log.Fatalf("Application failed to initialize: %v", err)
	}
	defer database.DBClose(app.DB)

	app.Run()
}
