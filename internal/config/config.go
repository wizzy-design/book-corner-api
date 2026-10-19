package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	UploadDir string
	AllowedOrigin string
}

func Load() (*Config, error) {
	err := godotenv.Load() // Loads the .env file

	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port: os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		UploadDir: os.Getenv("UPLOAD_DIR"),
		AllowedOrigin: os.Getenv("ALLOWED_ORIGIN"),
	}

	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return config, nil
}