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

func main() {
	fmt.Println("hello!")
	defer fmt.Println("good bye 👋")

	cfg := config.LoadConfig()
	log.Printf("config loaded: %+v\n", cfg)
	c := plex.New(cfg.Plex.Token, cfg.Plex.Address)
	vids, err := plex.GetRecentlyPlayed(c, "3")
	if err != nil {
		panic(err)
	}

	recentViews := make([]plex.VideoShort, 0, cfg.RecentMovieCount)
	log.Printf("making shortvid list of length %v\n", cfg.RecentMovieCount)
	for i, vid := range vids {
		if i > cfg.RecentMovieCount {
			break
		}

		recentViews = append(recentViews, plex.VideoShort{Title: vid.Title, Summary: vid.Summary, ContentRating: vid.ContentRating})
	}

	log.Printf("Recent views: %+v\n", recentViews)

	ollama, err := langchain.InitOllama(cfg.Ollama.Address, cfg.Ollama.Model)
	if err != nil {
		panic(err)
	}

	runSimple := os.Getenv("RUN_SIMPLE")
	full := runSimple == ""
	var recommendation string
	if full {

		recommendation, err = langchain.GenerateRecommendation(context.Background(), recentViews, &ollama)

	} else {
		recommendation, err = langchain.GenerateSimpleRecommendation(context.Background(), &ollama)
	}
	if err != nil {
		panic(err)
	}

	log.Printf("recommendation: %s\n", recommendation)
}
