package weaviate

import (
	"context"
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/schema"
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
	},
}

var CachedRecommendationClass = models.Class{
	Class:       cachedCollectionName,
	Description: "Schema for holding generated recommendations and the videos that generated them",
	Properties: []*models.Property{
		{
			Name:        "video_input",
			Description: "video_input is the latest watched videos provided as context to generate a recommendation",
			DataType:    []string{"text"},
		},
		{
			Name:        "generated_content",
			Description: "generated_content is the previously generated recommendated based on the provided recently watched videos",
			DataType:    []string{"text"},
		},
	},
}

func GetSchema() (*schema.Dump, error) {
	if err := createSchemaIfNotExists(&VideoClass); err != nil {
		return nil, err
	}

	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func createSchemaIfNotExists(class *models.Class) error {
	ok, err := client.Schema().ClassExistenceChecker().WithClassName(class.Class).Do(context.Background())
	if err != nil {
		log.Printf("could not check for class existence: %v\n", err)
	}

	if ok {
		log.Println("class exists, exiting")
		return nil
	}
	log.Println("class does not exist, creating")
	creator := client.Schema().ClassCreator().WithClass(class)
	if err := creator.Do(context.Background()); err != nil {
		return err
	}
	log.Println("created")

	return nil
}
