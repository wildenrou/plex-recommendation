package langchain

import (
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

// InitOllama is the entrypoint for interacting with Ollama provided
// LLMs
func InitOllama(address, languageModel, embeddingModel string) (*ollama.LLM, *ollama.LLM, error) {

	log.Println("initializing Ollama...")

	log.Println("initializing LLM...")
	var serverurl = "http://" + address + ":11434"
	llm, err := ollama.New(ollama.WithModel(languageModel), ollama.WithServerURL(serverurl))
	if err != nil {
		return nil, nil, err
	}
	log.Println("done")

	log.Println("initializng embedding model...")
	embeddingClient, err := ollama.New(
		ollama.WithModel(embeddingModel),
		ollama.WithServerURL(serverurl),
		ollama.WithKeepAlive("-1m"),
	)
	if err != nil {
		return nil, nil, err
	}
	log.Println("initialized")
	return llm, embeddingClient, nil
}
