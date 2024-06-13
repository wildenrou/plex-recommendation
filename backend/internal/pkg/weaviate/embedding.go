package weaviate

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

func embedChunkedDocument(ctx context.Context, embedder *ollama.LLM, texts []string) ([][]float32, error) {
	log.Println("start embed chunked documents")
	embeddings, err := embedder.CreateEmbedding(ctx, texts)
	if err != nil {
		return nil, err
	}
	log.Println("embed done!")
	return embeddings, nil
}
