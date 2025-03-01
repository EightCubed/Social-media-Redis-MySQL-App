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
)

// App encapsulates the application dependencies
type App struct {
	DB     *sql.DB
	Router *mux.Router
	Config config.Config
}

// Initialize sets up the database connection and router
func (a *App) Initialize() error {
	// Set up database connection
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

	// Verify connection
	err = a.DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Initialize router and routes
	a.Router = mux.NewRouter()
	a.initializeRoutes()

	return nil
}

// initializeRoutes sets up all API routes
func (a *App) initializeRoutes() {
	// apiRouter := a.Router.PathPrefix("/apis/v1").Subrouter()

	// Add your routes here
	// Example:
	// apiRouter.HandleFunc("/users", a.getUsers).Methods("GET")
}

// Run starts the web server
func (a *App) Run() {
	fmt.Printf("Starting server on :%s\n", a.Config.ServerPort)
	log.Fatal(http.ListenAndServe(":"+a.Config.ServerPort, a.Router))
}

// Close cleans up resources
func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}

func main() {
	// Load configuration from environment variables with defaults
	config := config.Config{
		DBHost:     getEnv("DB_HOST", "localhost:3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "rkn@1234"),
		DBName:     getEnv("DB_NAME", "mydatabase"), // Make sure to provide a default database name
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	// Initialize application
	app := App{Config: config}

	err := app.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Close()

	// Start the server
	app.Run()
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
