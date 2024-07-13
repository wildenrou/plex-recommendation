package weaviate

import (
	"context"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/codes"
	"log"

	"github.com/weaviate/weaviate/entities/models"
)

const (
	videoCollectionName  = "Videos"
	cachedCollectionName = "RecommendationsCache"
)

var VideoClass = models.Class{
	Class:       videoCollectionName,
	Description: "Schema for holding vectorized Plex video data",
	Properties: []*models.Property{
		{
			Name:        "title",
			Description: "title of the provided video",
			DataType:    []string{"text"},
		},
		{
			Name:        "summary",
			Description: "description of the video's plot",
			DataType:    []string{"text"},
		},
		{
			Name:        "content_rating",
			Description: "motion picture film association content rating",
			DataType:    []string{"text"},
		},
		{
			Name:        "plex_id",
			Description: "Plex GUID associated to the video",
			DataType:    []string{"text"},
		},
	},
}

func createSchemaIfNotExists(ctx context.Context, class *models.Class) error {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Create Schema If Not Exists"), telemetry.WithSpanPackage("weaviate"))
	defer span.End()
	ok, err := client.Schema().ClassExistenceChecker().WithClassName(class.Class).Do(context.Background())
	if err != nil {
		log.Printf("could not check for class existence: %v\n", err)
	}

	if ok {
		log.Println("class exists, exiting")
		span.SetStatus(codes.Ok, "class exists")
		return nil
	}
	log.Println("class does not exist, creating")
	creator := client.Schema().ClassCreator().WithClass(class)
	if err := creator.Do(context.Background()); err != nil {
		span.RecordError(err)
		return err
	}
	span.SetStatus(codes.Ok, "class created")
	log.Println("created")

	return nil
}
