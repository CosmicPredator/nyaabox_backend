package config

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

// LoadEnv loads the env vars from the .env file
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Error(".env file missing or of bad format")
		return err
	}
	return nil
}

// GetEnv gets the value of provided key else return error
func GetEnv(key string) (string, error) {
	s := os.Getenv(key)
	if s == "" {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return s, nil
}