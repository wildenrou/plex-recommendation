package langchain

import (
	"log"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// InitOllama is the entrypoint for interacting with Ollama provided
// LLMs
func InitOllama(address, model string) (llms.Model, *embeddings.EmbedderImpl, error) {

	log.Println("initializing Ollama...")
	var serverurl = "http://" + address + ":11434"
	llm, err := ollama.New(ollama.WithModel(model), ollama.WithServerURL(serverurl))
	if err != nil {
		return nil, nil, err
	}
	embedder, err := embeddings.NewEmbedder(llm, embeddings.WithBatchSize(512))
	if err != nil {
		return nil, nil, err
	}
	log.Println("initialized")
	return llm, embedder, nil
}
