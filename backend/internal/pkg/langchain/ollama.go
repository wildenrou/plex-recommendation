package langchain

import (
	"context"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/codes"
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

// InitOllama is the entrypoint for interacting with Ollama provided
// LLMs
func InitOllama(ctx context.Context, address, languageModel, embeddingModel string) (*ollama.LLM, *ollama.LLM, error) {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Ollama Initialization"), telemetry.WithSpanPackage("langchain"))
	defer span.End()
	log.Println("initializing Ollama...")

	log.Println("initializing LLM...")
	var serverurl = "http://" + address + ":11434"
	llm, err := ollama.New(ollama.WithModel(languageModel), ollama.WithServerURL(serverurl))
	if err != nil {
		span.RecordError(err)
		return nil, nil, err
	}
	span.AddEvent("initialized Language Model")
	log.Println("done")

	log.Println("initializng embedding model...")
	embeddingClient, err := ollama.New(
		ollama.WithModel(embeddingModel),
		ollama.WithServerURL(serverurl),
		ollama.WithKeepAlive("-1m"),
	)
	if err != nil {
		span.RecordError(err)
		return nil, nil, err
	}
	span.AddEvent("initialized Embedding Model")
	log.Println("initialized")
	span.SetStatus(codes.Ok, "initialized Ollama")
	return llm, embeddingClient, nil
}
