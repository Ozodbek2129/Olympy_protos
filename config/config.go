package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	USER_SERVICE string
	USER_ROUTER  string
	DB_HOST      string
	DB_PORT      string
	DB_PASSWORD  string
	DB_NAME      string
	SIGNING_KEY  string
	DB_USER      string
}

func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found?")
	}

	config := Config{}
	config.USER_SERVICE = cast.ToString(Coalesce("USER_SERVICE", ":50051"))
	config.USER_ROUTER = cast.ToString(Coalesce("USER_ROUTER", ":50052"))
	config.DB_USER=cast.ToString(Coalesce("DB_USER","postgres"))
	config.DB_HOST = cast.ToString(Coalesce("DB_HOST", "localhost"))
	config.DB_NAME = cast.ToString(Coalesce("DB_NAME", "olympy"))
	config.DB_PASSWORD = cast.ToString(Coalesce("DB_PASSWORD", "salom"))
	config.DB_PORT = cast.ToString(Coalesce("DB_PORT", "5432"))
	config.SIGNING_KEY = cast.ToString(Coalesce("SIGNING_KEY", "olympy"))

	return config
}

func Coalesce(key string, defaultValue interface{}) interface{} {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

