package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB_PORT string
	DB_USER string
	DB_PASS string
	DB_HOST string
	DB_NAME string

	SER_PORT string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := Config{
		DB_PORT:  os.Getenv("DB_PORT"),
		DB_USER:  os.Getenv("DB_USER"),
		DB_PASS:  os.Getenv("DB_PASS"),
		DB_HOST:  os.Getenv("DB_HOST"),
		DB_NAME:  os.Getenv("DB_NAME"),
		SER_PORT: os.Getenv("SER_PORT"),
	}
	return &cfg

}
