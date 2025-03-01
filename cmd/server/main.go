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
	connectionString := fmt.Sprintf("%s:%s@(%s)/%s",
		a.Config.DBUser,
		a.Config.DBPassword,
		a.Config.DBHost,
		a.Config.DBName)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}

	err = a.DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	err = database.CreateTables(a.DB)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	return nil
}

// initializeRoutes sets up all API routes
func (a *App) initializeRoutes() {
	apiRouter := a.Router.PathPrefix("/apis/v1").Subrouter()

	// Add your routes here
	// Example:
	apiRouter.HandleFunc("/health", a.healthCheckHandler).Methods("GET")
}

func (a *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (a *App) Run() {
	fmt.Printf("Starting server on :%s\n", a.Config.ServerPort)
	log.Fatal(http.ListenAndServe(":"+a.Config.ServerPort, a.Router))
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}

func main() {
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
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Close()

	app.Run()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
