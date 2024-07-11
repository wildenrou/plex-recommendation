package weaviate

import (
	"context"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/codes"
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

func embedChunkedDocument(ctx context.Context, embedder *ollama.LLM, texts []string) ([][]float32, error) {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Embed Chunked Document"), telemetry.WithSpanPackage("weaviate"))
	defer span.End()
	log.Println("start embed chunked documents")
	embeddings, err := embedder.CreateEmbedding(ctx, texts)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	log.Println("embed done!")
	span.SetStatus(codes.Ok, "embed done")
	return embeddings, nil
}
