package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Plex struct {
		Token                 string
		Address               string
		DefaultLibrarySection string
	}
	Ollama struct {
		Address        string
		LanguageModel  string
		EmbeddingModel string
	}
	Postgres struct {
		Host     string
		Username string
		Password string
		DBName   string
		Port     int
	}
	RecentMovieCount int
}

// loadEnv loads environment variables from a .env file.
func loadEnv() error {
	return godotenv.Load(".env")
}

// LoadConfig creates a Config struct based on current environment
func LoadConfig() *Config {
	// TODO: this is getting out of hand. Implement https://github.com/caarlos0/env
	if err := loadEnv(); err != nil {
		log.Printf("load env error: %s\n", err.Error())
	}

	var cfg Config
	if os.Getenv("PLEX_TOKEN") != "" {
		cfg.Plex.Token = os.Getenv("PLEX_TOKEN")
	}
	if os.Getenv("PLEX_ADDRESS") != "" {
		cfg.Plex.Address = os.Getenv("PLEX_ADDRESS")
	}

	// My movies library is at section 3, so I have the default set to that if
	// not provided via environment.
	cfg.Plex.DefaultLibrarySection = "3"
	if os.Getenv("PLEX_DEFAULT_LIBRARY_SECTION") != "" {
		cfg.Plex.DefaultLibrarySection = os.Getenv("PLEX_DEFAULT_LIBRARY_SECTION")
	}
	if os.Getenv("OLLAMA_ADDRESS") != "" {
		cfg.Ollama.Address = os.Getenv("OLLAMA_ADDRESS")
	}
	if os.Getenv("OLLAMA_LANGUAGE_MODEL") != "" {
		cfg.Ollama.LanguageModel = os.Getenv("OLLAMA_LANGUAGE_MODEL")
	}
	if os.Getenv("OLLAMA_EMBEDDING_MODEL") != "" {
		cfg.Ollama.EmbeddingModel = os.Getenv("OLLAMA_EMBEDDING_MODEL")
	}

	// Postgres values are defaulted to these initial values
	// but overriden by environment
	cfg.Postgres.Host = "postgres"
	if os.Getenv("POSTGRES_HOST") != "" {
		cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	}

	cfg.Postgres.Port = 5432
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Println("POSTGRES_PORT set but to no-int value")
	} else {
		cfg.Postgres.Port = port
	}

	cfg.Postgres.Username = "postgres"
	if os.Getenv("POSTGRES_USER") != "" {
		cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	}

	cfg.Postgres.Password = "postgres"
	if os.Getenv("POSTGRES_PASSWORD") != "" {
		cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	}

	cfg.Postgres.DBName = "caches"
	if os.Getenv("POSTGRES_DB") != "" {
		cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")
	}

	recentMovieCountStr := os.Getenv("RECENT_MOVIE_COUNT")
	count, err := strconv.Atoi(recentMovieCountStr)
	if recentMovieCountStr == "" || err != nil {
		cfg.RecentMovieCount = 5
	} else {
		cfg.RecentMovieCount = count
	}
	return &cfg
}
