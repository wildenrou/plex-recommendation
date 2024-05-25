package weaviate

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"

	"github.com/weaviate/weaviate/entities/models"
)

var client *weaviate.Client

func InitWeaviate() error {
	if client != nil {
		return nil
	}

	cfg := weaviate.Config{
		Host:   "weaviate:8080",
		Scheme: "http",
	}

	var err error
	client, err = weaviate.NewClient(cfg)
	if err != nil {
		return err
	}

	return nil
}

func InsertData(ctx context.Context, embedder embeddings.Embedder, videos []plex.VideoShort) error {
	log.Println("inserting data")
	defer log.Println("done!")

	var texts = make([]string, 0, len(videos))
	for _, video := range videos {
		texts = append(texts, video.String())
	}
	vectors, err := embedChunkedDocument(ctx, embedder, texts)
	if err != nil {
		return err
	}

	var objs = make([]*models.Object, 0, len(videos))
	for i, video := range videos {
		data := &models.Object{
			Class: videoCollectionName,
			Properties: map[string]any{
				"title":          video.Title,
				"summary":        video.Summary,
				"content_rating": video.ContentRating,
			},
			Vector: vectors[i],
		}

		objs = append(objs, data)
	}

	log.Println("start batch insert")
	defer log.Println("batch done!")
	batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objs...).Do(ctx)
	if err != nil {
		return err
	}

	var errors []string
	for _, res := range batchRes {
		if res.Result.Errors != nil {
			for _, err := range res.Result.Errors.Error {
				errors = append(errors, fmt.Sprintf("%v, ", err.Message))
			}

		}
	}

	if len(errors) != 0 {
		return fmt.Errorf("error in insert: %v", errors)
	}
	return nil
}

func QueryData(ctx context.Context, limit int) ([]*models.Object, error) {

	result, err := client.Data().ObjectsGetter().
		WithClassName(videoCollectionName).
		WithLimit(1).
		WithVector().
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return result, nil
}
