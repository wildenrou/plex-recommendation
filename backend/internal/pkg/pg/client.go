package pg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitPostgres() error {
	dsn := "host=postgres user=postgres password=postgres dbname=caches port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	db.AutoMigrate(&RecommendationCache{})

	return nil
}

func InsertData(input string, response string) error {
	cache := &RecommendationCache{
		InputTitles:     input,
		GeneratedOutput: response,
	}

	return db.Create(cache).Error
}

type queryOption struct {
	input    string
	response string
}

type QueryOption func(*queryOption)

func WithInputTitles(i string) QueryOption {
	return func(q *queryOption) {
		q.input = i
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

	var q RecommendationCache
	if query.input != "" {
		q.InputTitles = query.input
	}

	if query.response != "" {
		q.GeneratedOutput = query.response
	}

	var response *RecommendationCache
	result := db.Where(&q).First(response)
	return response, result.Error
}
