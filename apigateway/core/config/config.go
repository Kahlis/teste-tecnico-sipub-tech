package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env               string
	ListenPort        string
	GrpcServerAddress string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Env:               getEnv("ENV", "development"),
		ListenPort:        getEnv("LISTEN_PORT", "8080"),
		GrpcServerAddress: getEnv("GRPC_SERVER", "localhost:9090"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
