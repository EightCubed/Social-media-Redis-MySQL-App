package main

import (
	"fmt"
	"go-social-media/pkg/api/handlers"
	config "go-social-media/pkg/config"
	database "go-social-media/pkg/database"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type App struct {
	RedisClient *redis.Client
	DB          *database.DBConnection
	Router      *mux.Router
	Config      config.Config
}

func (a *App) startViewSync(interval time.Duration) {
	log.Printf("[INFO] Started Redis-to-MySQL view sync every %v", interval)
	go handlers.SyncViewsToDB(a.DB, a.RedisClient, interval)
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

	redisAddress := fmt.Sprintf("%s:%s", a.Config.RedisHost, a.Config.RedisPort)

	a.RedisClient = redis.NewClient(&redis.Options{
		Addr:         redisAddress,
		Password:     "",
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 5,
	})

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	log.Printf("Application initialization completed successfully.")
	return nil
}

func (a *App) initializeRoutes() {
	socialMediaHandler := handlers.ReturnHandler(a.DB, a.RedisClient)

	apiRouter := a.Router.PathPrefix("/apis/v1").Subrouter()

	// Health endpoints
	apiRouter.HandleFunc("/health", socialMediaHandler.HealthCheck).Methods("GET")

	// User endpoints
	apiRouter.HandleFunc("/user/{id:[0-9]+}", socialMediaHandler.GetUser).Methods("GET")
	apiRouter.HandleFunc("/user", socialMediaHandler.PostUser).Methods("POST")
	apiRouter.HandleFunc("/user", socialMediaHandler.ListUser).Methods("GET")
	apiRouter.HandleFunc("/user/{id:[0-9]+}", socialMediaHandler.UpdateUser).Methods("PATCH")
	apiRouter.HandleFunc("/user/{id:[0-9]+}", socialMediaHandler.DeleteUser).Methods("DELETE")

	// Post endpoints
	apiRouter.HandleFunc("/post/{id:[0-9]+}", socialMediaHandler.GetPost).Methods("GET")
	apiRouter.HandleFunc("/post", socialMediaHandler.PostPost).Methods("POST")
	apiRouter.HandleFunc("/post", socialMediaHandler.ListPost).Methods("GET")
	apiRouter.HandleFunc("/post/{id:[0-9]+}", socialMediaHandler.UpdatePost).Methods("PATCH")
	apiRouter.HandleFunc("/post/{id:[0-9]+}", socialMediaHandler.DeletePost).Methods("DELETE")

	// Metrics endpoint
	// a.Router.Handle("/metrics", promhttp.Handler())

	log.Printf("API routes initialized.")
}

func (a *App) Run() {
	srv := &http.Server{
		Addr:              ":" + a.Config.ServerPort,
		Handler:           a.Router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ReadHeaderTimeout: 2 * time.Second,
	}
	log.Printf("Starting server on port %s...\n", a.Config.ServerPort)
	err := srv.ListenAndServe()
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
		DBUser:      config.GetEnv("DB_USER", ""),
		DBPassword:  config.GetEnv("DB_PASSWORD", ""),
		DBName:      config.GetEnv("DB_NAME", "social_media_app"),
		ServerPort:  config.GetEnv("SERVER_PORT", "8080"),
		RedisHost:   config.GetEnv("REDIS_HOST", "redis-master.default.svc.cluster.local"),
		RedisPort:   config.GetEnv("REDIS_PORT", "6379"),
	}

	app := App{Config: config}

	err := app.Initialize()
	if err != nil {
		log.Fatalf("Application failed to initialize: %v", err)
	}
	defer database.DBClose(app.DB)

	app.startViewSync(handlers.CACHE_DURATION_LONG)

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil)) // Expose pprof endpoint
	}()

	app.Run()
}
