package main

import (
	"context"
	"fmt"
	"go-social-media/pkg/api/handlers"
	config "go-social-media/pkg/config"
	database "go-social-media/pkg/database"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type App struct {
	RedisClient *redis.Client
	DB          *database.DBConnection
	Router      *mux.Router
	Config      config.Config
	Server      *http.Server
	cancelFunc  context.CancelFunc
}

func (a *App) Initialize(ctx context.Context) error {
	log.Printf("Initializing application...")

	ctx, cancel := context.WithCancel(ctx)
	a.cancelFunc = cancel

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
		PoolSize:     500,
		MinIdleConns: 50,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  30 * time.Second,
		IdleTimeout:  5 * time.Minute,
	})

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	// Create HTTP Server with context support
	a.Server = &http.Server{
		Addr:              ":" + a.Config.ServerPort,
		Handler:           a.Router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ReadHeaderTimeout: 2 * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

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

	// Like endpoints
	apiRouter.HandleFunc("/post/{id:[0-9]+}/likes", socialMediaHandler.LikePost).Methods("POST")
	apiRouter.HandleFunc("/post/{id:[0-9]+}/likes", socialMediaHandler.LikeDelete).Methods("DELETE")

	// Metrics endpoint
	// a.Router.Handle("/metrics", promhttp.Handler())

	log.Printf("API routes initialized.")
}

func (a *App) Run(ctx context.Context) error {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on port %s...\n", a.Config.ServerPort)
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()

	go a.startViewSync(ctx, handlers.CACHE_DURATION_LONG)

	<-stopChan
	log.Println("Shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := a.Server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	if a.cancelFunc != nil {
		a.cancelFunc()
	}

	database.DBClose(a.DB)

	return nil
}

func (a *App) startViewSync(ctx context.Context, interval time.Duration) {
	log.Printf("[INFO] Started Redis-to-MySQL view sync every %v", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("View sync goroutine shutting down")
			return
		case <-ticker.C:
			handlers.SyncViewsToDB(a.DB, a.RedisClient, interval)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	err := app.Initialize(ctx)
	if err != nil {
		log.Fatalf("Application failed to initialize: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application run error: %v", err)
	}
}
