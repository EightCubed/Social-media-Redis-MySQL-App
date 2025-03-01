package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	config "go-social-media/pkg/config"
	database "go-social-media/pkg/database"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
	Config config.Config
}

func (a *App) Initialize() error {
	log.Println("Initializing application...")

	connectionString := fmt.Sprintf("%s:%s@(%s)/%s",
		a.Config.DBUser,
		a.Config.DBPassword,
		a.Config.DBHost,
		a.Config.DBName)

	var err error
	a.DB, err = database.DatabaseInit(connectionString, a.Config)
	if err != nil {
		return err
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	log.Println("Application initialization completed successfully.")
	return nil
}

func (a *App) initializeRoutes() {
	apiRouter := a.Router.PathPrefix("/apis/v1").Subrouter()
	apiRouter.HandleFunc("/health", a.healthCheckHandler).Methods("GET")

	log.Println("API routes initialized.")
}

func (a *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Health check endpoint hit.")

	err := a.DB.Ping()
	if err != nil {
		log.Printf("Health check failed: Database connection error: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(fmt.Sprintf("Database connection error: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (a *App) Run() {
	log.Printf("Starting server on port %s...\n", a.Config.ServerPort)
	err := http.ListenAndServe(":"+a.Config.ServerPort, a.Router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (a *App) Close() {
	if a.DB != nil {
		log.Println("Closing database connection.")
		a.DB.Close()
		log.Println("Database connection closed.")
	}
}

func main() {
	log.Println("Reading environment variables... 123")

	config := config.Config{
		DBHost:     getEnv("DB_HOST", "localhost:3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "rkn@1234"),
		DBName:     getEnv("DB_NAME", "mydatabase"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	app := App{Config: config}

	err := app.Initialize()
	if err != nil {
		log.Fatalf("Application failed to initialize: %v", err)
	}
	defer app.Close()

	app.Run()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Environment variable %s not set, using default: %s\n", key, defaultValue)
		return defaultValue
	}
	log.Printf("Using environment variable %s: %s\n", key, value)
	return value
}
