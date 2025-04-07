package handlers_test

import (
	"go-social-media/pkg/api/handlers"
	"go-social-media/pkg/database"
	"io"
	"log"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	log.SetOutput(io.Discard)
	RunSpecs(t, "Handlers Suite")
}

func createFakeDB() *database.DBConnection {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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

func createFakeSocialMediaHandler() *handlers.SocialMediaHandler {
	fakeDB := createFakeDB()
	fakeRedis := createFakeRedis()
	socialMediaHandler := handlers.ReturnHandler(fakeDB, fakeRedis)
	return socialMediaHandler
}
