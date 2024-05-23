package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PlexToken   string
	PlexAddress string
}

// loadEnv loads environment variables from a .env file.
func loadEnv() error {
	return godotenv.Load(".env")
}

// LoadConfig creates a Config struct based on current environment
func LoadConfig() *Config {
	if err := loadEnv(); err != nil {
		log.Printf("load env error: %s\n", err.Error())
	}

	var cfg Config
	if os.Getenv("PLEX_TOKEN") != "" {
		cfg.PlexToken = os.Getenv("PLEX_TOKEN")
	}
	if os.Getenv("PLEX_ADDRESS") != "" {
		cfg.PlexAddress = os.Getenv("PLEX_ADDRESS")
	}

	return &cfg
}
