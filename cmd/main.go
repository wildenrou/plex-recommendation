package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/config"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/langchain"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
)

const movieSection = "3"

func main() {
	fmt.Println("hello!")
	defer fmt.Println("good bye 👋")

	cfg := config.LoadConfig()
	c := plex.New(cfg.Plex.Token, cfg.Plex.Address)
	recentlyViewed, err := plex.GetRecentlyPlayed(c, movieSection, cfg.RecentMovieCount)
	if err != nil {
		panic(err)
	}

	fullCollection, err := plex.GetAllVideos(c, movieSection)
	if err != nil {
		panic(err)
	}

	ollama, err := langchain.InitOllama(cfg.Ollama.Address, cfg.Ollama.Model)
	if err != nil {
		panic(err)
	}

	runSimple := os.Getenv("RUN_SIMPLE")
	full := runSimple == ""
	var recommendation string
	if full {

		recommendation, err = langchain.GenerateRecommendation(context.Background(), recentlyViewed, fullCollection, &ollama)

	} else {
		recommendation, err = langchain.GenerateSimpleRecommendation(context.Background(), &ollama)
	}
	if err != nil {
		panic(err)
	}

	log.Printf("%s\n", recommendation)
}
