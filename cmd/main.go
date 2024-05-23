package main

import (
	"fmt"
	"log"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/config"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
)

func main() {
	fmt.Println("hello!")
	defer fmt.Println("good bye 👋")

	cfg := config.LoadConfig()

	c := plex.New(cfg.PlexToken, cfg.PlexAddress)
	vids, err := plex.GetRecentlyPlayed(c, "3")
	if err != nil {
		panic(err)
	}

	for _, video := range vids {
		log.Printf("title: %s\ndescription: %s\nrating: %s\n", video.Title, video.Summary, video.ContentRating)
	}
}
