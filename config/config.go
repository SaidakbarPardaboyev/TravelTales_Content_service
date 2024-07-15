package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	AUTH_SERVICE_PORT    string
	USER_SERVICE_PORT    string
	CONTENT_SERVICE_PORT string
	API_GATEWAY_PORT     string
	DB_HOST              string
	DB_PORT              string
	DB_USER              string
	DB_NAME              string
	DB_PASSWORD          string
	SINGNING_KEY_ACCESS  string
	SINGNING_KEY_REFRESH string
	EMAIL                string
	PASSWORD             string
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf(".env file not found: %s", err)
	}
	config := Config{}
	config.AUTH_SERVICE_PORT = cast.ToString(coalesce("AUTH_SERVICE_PORT", ":8080"))
	config.USER_SERVICE_PORT = cast.ToString(coalesce("USER_SERVICE_PORT", ":8080"))
	config.CONTENT_SERVICE_PORT = cast.ToString(coalesce("CONTENT_SERVICE_PORT", ":8080"))
	config.API_GATEWAY_PORT = cast.ToString(coalesce("API_GATEWAY_PORT", ":8080"))
	config.DB_HOST = cast.ToString(coalesce("DB_HOST", ":8080"))
	config.DB_PORT = cast.ToString(coalesce("DB_PORT", ":8080"))
	config.DB_USER = cast.ToString(coalesce("DB_USER", ":8080"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", ":8080"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", ":8080"))
	config.SINGNING_KEY_ACCESS = cast.ToString(coalesce("SINGNING_KEY_ACCESS", ":8080"))
	config.SINGNING_KEY_REFRESH = cast.ToString(coalesce("SINGNING_KEY_REFRESH", ":8080"))
	config.EMAIL = cast.ToString(coalesce("EMAIL", ":8080"))
	config.PASSWORD = cast.ToString(coalesce("PASSWORD", ":8080"))

	return &config
}

func coalesce(key string, defaultValue interface{}) interface{} {
	if res, exists := os.LookupEnv(key); exists {
		return res
	}
	return defaultValue
}
