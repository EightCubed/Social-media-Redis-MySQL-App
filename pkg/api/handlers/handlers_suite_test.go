package handlers_test

import (
	"go-social-media/pkg/api/handlers"
	"go-social-media/pkg/database"
	"io"
	"log"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	log.SetOutput(io.Discard)
	RunSpecs(t, "Handlers Suite")
}

func createFakeDB() *database.DBConnection {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	database.AutoMigrateTables(db)

	return &database.DBConnection{
		GormDBWriter: db,
		GormDBReader: db,
	}
}

func createFakeRedis() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

func createFakeSocialMediaHandlerAndRouter() (*handlers.SocialMediaHandler, *mux.Router) {
	fakeDB := createFakeDB()
	fakeRedis := createFakeRedis()
	socialMediaHandler := handlers.ReturnHandler(fakeDB, fakeRedis)

	router := mux.NewRouter()
	registerRoutes(router, socialMediaHandler)
	return socialMediaHandler, router
}

func registerRoutes(router *mux.Router, socialMediaHandler *handlers.SocialMediaHandler) {
	apiRouter := router.PathPrefix("/apis/v1").Subrouter()

	// Health endpoints
	apiRouter.HandleFunc("/health", socialMediaHandler.HealthCheck).Methods("GET")

	// User endpoints
	apiRouter.HandleFunc("/user", socialMediaHandler.ListUser).Methods("GET")
	apiRouter.HandleFunc("/user", socialMediaHandler.PostUser).Methods("POST")
	apiRouter.HandleFunc("/user/{user_id:[0-9]+}", socialMediaHandler.GetUser).Methods("GET")
	apiRouter.HandleFunc("/user/{user_id:[0-9]+}", socialMediaHandler.UpdateUser).Methods("PATCH")
	apiRouter.HandleFunc("/user/{user_id:[0-9]+}", socialMediaHandler.DeleteUser).Methods("DELETE")

	// Post endpoints
	apiRouter.HandleFunc("/post", socialMediaHandler.ListPost).Methods("GET")
	apiRouter.HandleFunc("/post", socialMediaHandler.PostPost).Methods("POST")
	apiRouter.HandleFunc("/post/{post_id:[0-9]+}", socialMediaHandler.GetPost).Methods("GET")
	apiRouter.HandleFunc("/post/{post_id:[0-9]+}", socialMediaHandler.UpdatePost).Methods("PATCH")
	apiRouter.HandleFunc("/post/{post_id:[0-9]+}", socialMediaHandler.DeletePost).Methods("DELETE")

	// Like endpoints
	apiRouter.HandleFunc("/post/{post_id:[0-9]+}/likes", socialMediaHandler.LikePost).Methods("POST")
	apiRouter.HandleFunc("/post/{post_id:[0-9]+}/likes", socialMediaHandler.LikeDelete).Methods("DELETE")

	// Comment endpoints
	apiRouter.HandleFunc("/post/{post_id:[0-9]+}/comments", socialMediaHandler.ListComments).Methods("GET")
	apiRouter.HandleFunc("/post/{post_id:[0-9]+}/comments", socialMediaHandler.PostComment).Methods("POST")
	apiRouter.HandleFunc("/comments/{comment_id:[0-9]+}", socialMediaHandler.UpdateComment).Methods("PATCH")
	apiRouter.HandleFunc("/comments/{comment_id:[0-9]+}", socialMediaHandler.DeleteComment).Methods("DELETE")
}
