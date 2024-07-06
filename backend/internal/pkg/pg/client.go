package pg

import (
	b64 "encoding/base64"
	"log"
	"slices"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var client *gorm.DB

func InitPostgres() error {
	if client != nil {
		return nil
	}
	dsn := "host=postgres user=postgres password=postgres dbname=caches port=5432 sslmode=disable"
	var err error
	client, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	log.Println("automigrating db")
	if err := client.AutoMigrate(&RecommendationCache{}); err != nil {
		return err
	}

	return nil
}

func InsertData(input []string, response string) error {
	// sort the incoming titles slice so recently viewed is
	// indifferent to order of recent viewing.
	slices.Sort(input)
	cache := &RecommendationCache{
		InputTitles:     toBase64(buildStringFromSlice(input)),
		GeneratedOutput: response,
	}

	return client.Create(cache).Error
}

type queryOption struct {
	input    string
	response string
}

type QueryOption func(*queryOption)

func WithInputTitles(i string) QueryOption {
	return func(q *queryOption) {
		q.input = b64.StdEncoding.EncodeToString([]byte(i))
	}
}

func WithReponse(r string) QueryOption {
	return func(q *queryOption) {
		q.response = r
	}
}

func QueryData(opts ...QueryOption) (*RecommendationCache, error) {
	var query = &queryOption{}
	for _, opt := range opts {
		opt(query)
	}

	var q = RecommendationCache{}
	if query.input != "" {
		q.InputTitles = query.input
	}

	if query.response != "" {
		q.GeneratedOutput = query.response
	}
	var response = RecommendationCache{}
	result := client.Where(&q).First(&response)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return &response, nil
}
