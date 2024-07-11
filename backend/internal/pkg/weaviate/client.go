package weaviate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"

	"github.com/weaviate/weaviate/entities/models"
)

var client *weaviate.Client

type queryOption struct {
	className string
	limit     int
}

type QueryOption func(*queryOption)

func WithClassName(s string) QueryOption {
	return func(q *queryOption) {
		q.className = s
	}
}

func WithLimit(i int) QueryOption {
	return func(q *queryOption) {
		q.limit = i
	}
}

type insertOption struct {
	videos []plex.VideoShort
}

type InsertOption func(*insertOption)

func WithVideos(v []plex.VideoShort) InsertOption {
	return func(i *insertOption) {
		i.videos = v
	}
}

func InitWeaviate(ctx context.Context, c plex.Client, embedder *ollama.LLM) error {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("InitWeaviate"))
	defer span.End()
	if client != nil {
		span.SetStatus(codes.Ok, "Connected to Weaviate Previously")
		return nil
	}

	cfg := weaviate.Config{
		Host:   "weaviate:8080",
		Scheme: "http",
	}

	var err error
	client, err = weaviate.NewClient(cfg)
	if err != nil {
		span.RecordError(err)
		return err
	}

	classesToCheck := []models.Class{VideoClass}

	for _, class := range classesToCheck {
		if err := createSchemaIfNotExists(&class); err != nil {
			span.RecordError(err)
			return err
		}
	}

	if err := insertPlexMedia(ctx, c, embedder); err != nil {
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "Connected to Weaviate")
	return nil
}

func InsertData(ctx context.Context, embedder *ollama.LLM, opts ...InsertOption) error {
	log.Println("inserting data")
	defer log.Println("done!")

	options := &insertOption{}
	for _, opt := range opts {
		opt(options)
	}

	var objs = make([]*models.Object, 0)
	if options.videos != nil {
		var texts = make([]string, 0, len(options.videos))
		for _, video := range options.videos {
			texts = append(texts, video.String())
		}
		vectors, err := embedChunkedDocument(ctx, embedder, texts)
		if err != nil {
			return err
		}

		for i, video := range options.videos {
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

func QueryData(ctx context.Context, opts ...QueryOption) ([]*models.Object, error) {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Query Data"))
	defer span.End()
	span.SetAttributes(attribute.String("package", "weaviate"))

	options := &queryOption{}
	for _, opt := range opts {
		opt(options)
	}

	limit := 5
	if options.className == "" {
		err := errors.New("no class provided to required WithClassName option")
		span.RecordError(err)
		return nil, err
	}

	if options.limit > 0 {
		limit = options.limit
	}

	allObjects := make([]*models.Object, 0)
	after := ""
	for {
		getter := client.Data().ObjectsGetter().
			WithClassName(options.className).
			WithLimit(limit).
			WithVector()

		if after != "" {
			getter = getter.WithAfter(after)
		}

		result, err := getter.Do(ctx)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		allObjects = append(allObjects, result...)
		if len(result) <= limit {
			break
		}
		after = result[len(result)-1].ID.String()
	}

	span.SetStatus(codes.Ok, "query complete")
	return allObjects, nil
}

func insertPlexMedia(ctx context.Context, c plex.Client, embedder *ollama.LLM) error {
	log.Println("performing migration on load...")
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Insert Plex Media"))
	defer span.End()
	vids, err := plex.GetAllVideos(ctx, c, "3")
	if err != nil {
		span.RecordError(err)
		return err
	}
	log.Println("got ", len(vids), " videos")
	span.SetAttributes(attribute.Int("count", len(vids)))
	savedData, err := QueryData(ctx, WithClassName(VideoClass.Class), WithLimit(500))
	if err != nil {
		span.RecordError(err)
		return err
	}
	span.AddEvent("saved data")

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
		if err = InsertData(ctx, embedder, WithVideos(toSave)); err != nil {
			span.RecordError(err)
			return err
		}

		span.AddEvent("saved found diff data")
	}

	span.SetStatus(codes.Ok, "migration complete")

	log.Println("complete")
	return nil
}

func VectorQuery(ctx context.Context, collectionName string, limit int, vectors [][]float32) ([]*plex.VideoShort, error) {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Vector Query"))
	defer span.End()
	span.SetAttributes(attribute.String("package", "weaviate"))
	nearVectorArgument := client.GraphQL().NearVectorArgBuilder()
	for _, vector := range vectors {
		nearVectorArgument.WithVector(vector)
	}
	fields := []graphql.Field{
		{Name: "title"},
		{Name: "summary"},
		{Name: "content_rating"},
	}
	resp, err := client.GraphQL().Get().WithClassName(collectionName).WithFields(fields...).WithNearVector(nearVectorArgument).WithLimit(limit).Do(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.AddEvent("query successful")

	if resp.Errors != nil {
		var errs string
		for _, err := range resp.Errors {
			errs += err.Message + "\n"
		}
		span.RecordError(errors.New(errs))
		return nil, errors.New(errs)
	}

	results, err := resp.MarshalBinary()
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.AddEvent("marshall binary successful")

	type marshalResults struct {
		Data struct {
			Get struct {
				Videos []*plex.VideoShort `json:"Videos"`
			} `json:"Get"`
		} `json:"data"`
	}

	var toReturn marshalResults
	if err := json.Unmarshal(results, &toReturn); err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "query successful")

	return toReturn.Data.Get.Videos, nil
}
