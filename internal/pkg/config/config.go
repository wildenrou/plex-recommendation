package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Plex struct {
		Token   string
		Address string
	}
	Ollama struct {
		Address        string
		LanguageModel  string
		EmbeddingModel string
	}
	RecentMovieCount int
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
		cfg.Plex.Token = os.Getenv("PLEX_TOKEN")
	}
	if os.Getenv("PLEX_ADDRESS") != "" {
		cfg.Plex.Address = os.Getenv("PLEX_ADDRESS")
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

	recentMovieCountStr := os.Getenv("RECENT_MOVIE_COUNT")
	count, err := strconv.Atoi(recentMovieCountStr)
	if recentMovieCountStr == "" || err != nil {
		cfg.RecentMovieCount = 5
	} else {
		cfg.RecentMovieCount = count
	}
	return &cfg
}
