package config

import (
	"os"
)

type Config struct {
	DBWriteHost string
	DBReadHost  string
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
	RedisHost   string
	RedisPort   string
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
