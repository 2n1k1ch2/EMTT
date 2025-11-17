package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DB_PORT  string
	DB_USER  string
	DB_PASS  string
	DB_HOST  string
	DB_NAME  string
	SER_PORT string
}

func Load() *Config {

	err := godotenv.Load(".env")
	log.Printf("Failed to load .env file: %s", err.Error())
	cfg := Config{
		DB_PORT:  getEnv("DB_PORT", "5432"),
		DB_USER:  getEnv("DB_USER", "postgres"),
		DB_PASS:  getEnv("DB_PASS", "password"),
		DB_HOST:  getEnv("DB_HOST", "postgres"),
		DB_NAME:  getEnv("DB_NAME", "postgres"),
		SER_PORT: getEnv("SER_PORT", "8088"),
	}

	log.Printf("Loaded config: DB_HOST=%s, DB_PORT=%s, SER_PORT=%s",
		cfg.DB_HOST, cfg.DB_PORT, cfg.SER_PORT)

	return &cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
