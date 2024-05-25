package weaviate

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/embeddings"
)

func embedChunkedDocument(ctx context.Context, embedder embeddings.Embedder, texts []string) ([][]float32, error) {
	log.Println("start embed chunked documents")
	embeddings, err := embedder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, err
	}
	log.Println("embed done!")
	return embeddings, nil
}
