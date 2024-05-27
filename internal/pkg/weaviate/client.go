package weaviate

import (
	"context"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"

	"github.com/weaviate/weaviate/entities/models"
)

var client *weaviate.Client

func InitWeaviate(c plex.Client, embedder embeddings.Embedder) error {
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

	if err := createSchemaIfNotExists(); err != nil {
		return err
	}

	if err := insertPlexMedia(c, embedder); err != nil {
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
	allObjects := make([]*models.Object, 0)
	after := ""
	for {
		getter := client.Data().ObjectsGetter().
			WithClassName(videoCollectionName).
			WithVector().
			WithLimit(limit)

		if after != "" {
			getter = getter.WithAfter(after)
		}

		result, err := getter.Do(ctx)
		if err != nil {
			return nil, err
		}
		allObjects = append(allObjects, result...)
		if len(result) <= limit {
			break
		}
		after = result[len(result)-1].ID.String()
	}

	return allObjects, nil
}

func insertPlexMedia(c plex.Client, embedder embeddings.Embedder) error {
	log.Println("performing migration on load...")
	vids, err := plex.GetAllVideos(c, "3")
	if err != nil {
		return err
	}
	log.Println("got ", len(vids), " videos")

	savedData, err := QueryData(context.Background(), 500)
	if err != nil {
		return err
	}

	log.Println("found ", len(savedData), " videos in the db")

	// map for faster lookup when we check if a video is
	// already saved
	savedHm := make(map[string]strfmt.UUID, len(savedData))
	for _, obj := range savedData {
		summary := obj.Properties.(map[string]interface{})["summary"].(string)
		savedHm[summary] = obj.ID
	}

	toSave := make([]plex.VideoShort, 0, len(vids))
	for _, vid := range vids {
		if _, ok := savedHm[vid.Summary]; !ok {
			// this video not found in the saved video
			// map, so add it to the list of new media
			// to save
			toSave = append(toSave, vid)
		}
	}

	log.Println("found ", len(toSave), " videos to save")
	if len(toSave) > 0 {
		if err = InsertData(context.Background(), embedder, toSave); err != nil {
			return err
		}
	}

	log.Println("complete")
	return nil
}
