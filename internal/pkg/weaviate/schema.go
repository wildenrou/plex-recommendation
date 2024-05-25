package weaviate

import (
	"context"
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/schema"
	"github.com/weaviate/weaviate/entities/models"
)

const (
	videoCollectionName = "videos"
)

var VideoClass = models.Class{
	Class:       videoCollectionName,
	Description: "Schema for holding vectorized Plex video data",
	Properties: []*models.Property{
		{
			Name:        "title",
			Description: "title of the provided data file",
			DataType:    []string{"text"},
		},
		{
			Name:        "summary",
			Description: "string content of the data",
			DataType:    []string{"text[]"},
		},
		{
			Name:        "content_rating",
			Description: "motion picture film association content rating",
			DataType:    []string{"text"},
		},
	},
}

func GetSchema() (*schema.Dump, error) {
	if err := CreateSchemaIfNotExists(); err != nil {
		return nil, err
	}

	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func CreateSchemaIfNotExists() error {
	ok, err := client.Schema().ClassExistenceChecker().WithClassName(videoCollectionName).Do(context.Background())
	if err != nil {
		log.Printf("could not check for class existence: %v\n", err)
	}

	if ok {
		log.Println("class exists, exiting")
		return nil
	}
	log.Println("class does not exist, creating")
	creator := client.Schema().ClassCreator().WithClass(&VideoClass)
	if err := creator.Do(context.Background()); err != nil {
		return err
	}
	log.Println("created")

	return nil
}
