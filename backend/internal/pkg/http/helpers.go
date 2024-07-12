package httpinternal

import (
	"context"
	"fmt"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/codes"
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
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Get Recommendation"))
	defer span.End()
	recentlyViewed, err := plex.GetRecentlyPlayed(ctx, plexClient, section, limit)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
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
	resp, err := pg.QueryData(ctx, pg.WithInputTitles(titles))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		log.Println("could not query cache for these titles: ", err.Error())
	}

	if resp.GeneratedOutput != "" {
		log.Println("found cached recommendation")
		span.SetStatus(codes.Ok, "found cached recommendation")
		span.AddEvent("cache found")
		return resp.GeneratedOutput, nil
	}

	span.AddEvent("no cached recommendation")

	log.Println("embeding recently viewed...")
	log.Println("embedding ", len(rvTexts), " texts")
	rvEmbeddings, err := ollamaEmbedder.CreateEmbedding(ctx, rvTexts)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	span.AddEvent("embeddings complete")
	log.Println("embeddings complete, querying database")

	results, err := weaviate.VectorQuery(ctx, weaviate.VideoClass.Class, limit, rvEmbeddings)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	span.AddEvent("vector query complete")
	log.Println("complete")

	rvStr := buildStringFromSlice(results)

	fullCollection, err := plex.GetAllVideos(ctx, plexClient, section)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}

	fcStr := buildStringFromSlice(fullCollection)

	recommendation, err := langchain.GenerateRecommendation(ctx, rvStr, fcStr, ollamaLlm)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	span.AddEvent("recommend complete")
	normalized, err := langchain.NormalizeLLMResponse(ctx, recommendation, ollamaLlm)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	span.AddEvent("normalization complete")
	// save this generated text back to the db
	if err := pg.InsertData(ctx, titles, normalized); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.AddEvent("insert failed")
		log.Println("could not cache this response: ", err.Error())
	}
	span.SetStatus(codes.Ok, "generation completed")
	return normalized, nil

}
