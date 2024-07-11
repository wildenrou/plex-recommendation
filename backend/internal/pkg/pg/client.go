package pg

import (
	"context"
	b64 "encoding/base64"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func InsertData(ctx context.Context, input []string, response string) error {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("InsertData"))
	defer span.End()
	span.SetAttributes(attribute.String("package", "pg"))
	// sort the incoming titles slice so recently viewed is
	// indifferent to order of recent viewing.
	slices.Sort(input)
	cache := &RecommendationCache{
		InputTitles:     toBase64(buildStringFromSlice(input)),
		GeneratedOutput: response,
	}
	if err := client.Create(cache).Error; err != nil {
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "insert complete")
	return nil
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

func QueryData(ctx context.Context, opts ...QueryOption) (*RecommendationCache, error) {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("QueryData"))
	defer span.End()

	span.SetAttributes(attribute.String("package", "pg"))
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
		span.RecordError(result.Error)
		return nil, result.Error
	}
	span.SetStatus(codes.Ok, "query succeeded")
	return &response, nil
}
