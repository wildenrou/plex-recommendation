package httpinternal

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/langchain"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
)

func getRecommendation(w http.ResponseWriter, r *http.Request) {
	section := r.PathValue("movieSection")
	var limit int
	limitQuery, ok := r.URL.Query()["limit"]
	if ok {
		limit, _ = strconv.Atoi(limitQuery[0])
	}

	recentlyViewed, err := plex.GetRecentlyPlayed(plexClient, section, limit)
	if err != nil {
		panic(err)
	}

	rvStr := buildStringFromSlice(recentlyViewed)

	fullCollection, err := plex.GetAllVideos(plexClient, section)
	if err != nil {
		panic(err)
	}

	fcStr := buildStringFromSlice(fullCollection)

	runSimple := os.Getenv("RUN_SIMPLE")
	full := runSimple == ""
	var recommendation string
	if full {

		recommendation, err = langchain.GenerateRecommendation(context.Background(), rvStr, fcStr, &ollamaLlm)

	} else {
		recommendation, err = langchain.GenerateSimpleRecommendation(context.Background(), &ollamaLlm)
	}
	if err != nil {
		panic(err)
	}

	log.Printf("\n%s\n", recommendation)
}
