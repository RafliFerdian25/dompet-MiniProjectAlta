package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIPort      string
	DB_ADDRESS   string
	DB_USERNAME  string
	DB_PASSWORD  string
	DB_NAME      string
	TOKEN_SECRET string
}

// Config func to get env value from key ---
func ConfigValue(key string) string{
    // load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    return os.Getenv(key)
}