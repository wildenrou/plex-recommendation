package httpinternal

import (
	"context"
	"fmt"
	"log"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/langchain"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/pg"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/weaviate"
)

func buildStringFromSlice[T any](slice []T) string {
	return fmt.Sprintf("%+v", slice)
}

func getRecommendation(ctx context.Context, section string, limit int) (string, error) {
	recentlyViewed, err := plex.GetRecentlyPlayed(plexClient, section, limit)
	if err != nil {
		return "", err
	}

	// LLM inputs operate on strings, so force the structs from the call to
	// plex into their stringified forms
	rvTexts := make([]string, 0, len(recentlyViewed))
	titles := make([]string, 0, len(recentlyViewed))
	for _, vid := range recentlyViewed {
		rvTexts = append(rvTexts, vid.String())
		titles = append(titles, vid.Title)
	}

	// query the cache to see if we've asked for recommendations
	// based on this exact recently viewed
	resp, err := pg.QueryData(pg.WithInputTitles(buildStringFromSlice(titles)))
	if err != nil {
		log.Println("could not query cache for these titles: ", err.Error())
	}

	if resp.GeneratedOutput != "" {
		log.Println("found cached recommendation")
		return resp.GeneratedOutput, nil
	}

	log.Println("embeding recently viewed...")
	log.Println("embedding ", len(rvTexts), " texts")
	rvEmbeddings, err := ollamaEmbedder.CreateEmbedding(ctx, rvTexts)
	if err != nil {
		return "", err
	}

	log.Println("embeddings complete, querying database")

	results, err := weaviate.VectorQuery(context.Background(), weaviate.VideoClass.Class, limit, rvEmbeddings)
	if err != nil {
		return "", err
	}

	log.Println("complete")

	rvStr := buildStringFromSlice(results)

	fullCollection, err := plex.GetAllVideos(plexClient, section)
	if err != nil {
		return "", err
	}

	fcStr := buildStringFromSlice(fullCollection)

	recommendation, err := langchain.GenerateRecommendation(ctx, rvStr, fcStr, ollamaLlm)
	if err != nil {
		return "", err
	}

	normalized, err := langchain.NormalizeLLMResponse(ctx, recommendation, ollamaLlm)
	if err != nil {
		return "", err
	}

	// save this generated text back to the db
	if err := pg.InsertData(titles, normalized); err != nil {
		log.Println("could not cache this response: ", err.Error())
	}

	return normalized, nil

}
