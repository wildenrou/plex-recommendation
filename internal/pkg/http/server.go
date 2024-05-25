package httpinternal

import (
	"log"
	"net/http"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/config"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/langchain"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/weaviate"
)

var (
	plexClient     *plex.PlexClient
	ollamaLlm      llms.Model
	ollamaEmbedder embeddings.Embedder
)

// StartServer initializes dependent services that are
// required to handle HTTP requests. This is blocking.
func StartServer(c *config.Config, shutdownChan chan error) {
	initLLM(c)
	initPlex(c)
	initVectorStore()
	initHttpServer(shutdownChan)
}

// initHttpServer is a blocking function that runs the HTTP server
// and passes requests to their respective handler functions
func initHttpServer(s chan error) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /recommendation/{movieSection}", getRecommendation)
	log.Println("serving http...")
	if err := http.ListenAndServe(":8090", mux); err != nil {
		s <- err
	}
	s <- nil
}

// initPlex creates the Plex client the server uses to
// connect to and execute queries against Plex
func initPlex(c *config.Config) {
	log.Println("initialzing plex client...")
	defer log.Println("initialized")
	if plexClient != nil {
		return
	}
	plexClient = plex.New(c.Plex.Token, c.Plex.Address)
}

// initLLM creates the Ollama LLM client the server uses
// to connect to and execute generation and embeddings
func initLLM(c *config.Config) error {
	if ollamaLlm != nil {
		return nil
	}
	var err error
	ollamaLlm, ollamaEmbedder, err = langchain.InitOllama(c.Ollama.Address, c.Ollama.Model)
	if err != nil {
		return err
	}
	return nil
}

// initVectorStore connects to Weaviate for storing
// Plex data and related embeddings
func initVectorStore() error {
	if err := weaviate.InitWeaviate(); err != nil {
		return err
	}
	return nil
}
