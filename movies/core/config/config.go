package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env          string
	DbName       string
	DBUri        string
	DBUser       string
	DbCollection string
	DBPassword   string
	ListenPort   string
	ApiPort      string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Env:          getEnv("ENV", "development"),
		DbName:       getEnv("MONGO_DB", "mongo"),
		DBUser:       getEnv("MONGO_DB_USER", "admin"),
		DBPassword:   getEnv("MONGO_DB_PASSWORD", "password"),
		DbCollection: getEnv("MONGO_DB_COLLECTION", "movies"),
		DBUri:        getEnv("MONGO_DB_URI", "mongodb://admin:password@mongodb:27017/movies?authSource=admin"),
		ListenPort:   getEnv("LISTEN_PORT", "9090"),
		ApiPort:      getEnv("API_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
