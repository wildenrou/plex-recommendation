package httpinternal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/langchain"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/pg"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/weaviate"
)

func formatHttpError(err error) []byte {
	return []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
}

func getRecommendation(w http.ResponseWriter, r *http.Request) {
	section := r.PathValue("movieSection")
	var limit int
	limitQuery, ok := r.URL.Query()["limit"]
	if ok {
		limit, _ = strconv.Atoi(limitQuery[0])
	}

	recentlyViewed, err := plex.GetRecentlyPlayed(plexClient, section, limit)
	if err != nil {
		w.Write(formatHttpError(err))
		return
	}

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

	var normalized string

	// TODO: refactor if/else logic to functions
	if resp.GeneratedOutput != "" {
		normalized = resp.GeneratedOutput
		log.Println("found cached recommendation")
	} else {
		log.Println("embeding recently viewed...")
		log.Println("embedding ", len(rvTexts), " texts")
		rvEmbeddings, err := ollamaEmbedder.CreateEmbedding(r.Context(), rvTexts)
		if err != nil {
			w.Write(formatHttpError(err))
			return
		}

		log.Println("embeddings complete, querying database")

		results, err := weaviate.VectorQuery(context.Background(), weaviate.VideoClass.Class, limit, rvEmbeddings)
		if err != nil {
			w.Write(formatHttpError(err))
			return
		}

		log.Println("complete")

		rvStr := buildStringFromSlice(results)

		fullCollection, err := plex.GetAllVideos(plexClient, section)
		if err != nil {
			w.Write(formatHttpError(err))
			return
		}

		fcStr := buildStringFromSlice(fullCollection)

		runSimple := os.Getenv("RUN_SIMPLE")
		full := runSimple == ""
		var recommendation string
		if full {
			recommendation, err = langchain.GenerateRecommendation(context.Background(), rvStr, fcStr, ollamaLlm)

		} else {
			recommendation, err = langchain.GenerateSimpleRecommendation(context.Background(), ollamaLlm)
		}
		if err != nil {
			w.Write(formatHttpError(err))
			return
		}

		normalized, err = langchain.NormalizeLLMResponse(r.Context(), recommendation, ollamaLlm)
		if err != nil {
			w.Write(formatHttpError(err))
			return
		}

		// save this generated text back to the db
		if err := pg.InsertData(buildStringFromSlice(titles), normalized); err != nil {
			log.Println("could not cache this response: ", err.Error())
		}
	}

	var respStruct []*plex.VideoShort
	if err := json.Unmarshal([]byte(normalized), &respStruct); err != nil {
		w.Write(formatHttpError(err))
		return
	}
	respBytes, err := json.Marshal(&respStruct)
	if err != nil {
		w.Write(formatHttpError(err))
		return
	}

	_, err = w.Write(respBytes)
	if err != nil {
		log.Println("could not write back to client: ", err.Error())
	}
}
