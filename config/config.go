package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	API_PORT      string
	DB_ADDRESS   string
	DB_USERNAME  string
	DB_PASSWORD  string
	DB_NAME      string
	TOKEN_SECRET string
}

var Cfg *Config

func InitConfig() {
	cfg := &Config{}

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
	
	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Println(err)
	}
	

	Cfg = cfg
}

// Config func to get env value from key ---
func ConfigValue(key string) string{
    // load .env file
    err := godotenv.Load(".env")
    if err != nil {
        fmt.Print("Error loading .env file", err)
    }
	fmt.Println(os.Getenv("API_PORT"))
    return os.Getenv(key)
}