package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Environment variable %s not set, using default: %s\n", key, defaultValue)
		return defaultValue
	}
	log.Printf("Using environment variable %s: %s\n", key, value)
	return value
}
