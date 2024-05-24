package httpinternal

import (
	"net/http"

	"github.com/tmc/langchaingo/llms"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/config"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/langchain"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
)

var (
	plexClient *plex.PlexClient
	ollamaLlm  llms.Model
)

// StartServer initializes dependent services that are
// required to handle HTTP requests. This is blocking.
func StartServer(c *config.Config, shutdownChan chan error) {
	initLLM(c)
	initPlex(c)
	initHttpServer(shutdownChan)
}

// initHttpServer is a blocking function that runs the HTTP server
// and passes requests to their respective handler functions
func initHttpServer(s chan error) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /recommendation/{movieSection}", getRecommendation)
	if err := http.ListenAndServe("localhost:8090", mux); err != nil {
		s <- err
	}
	s <- nil
}

// initPlex creates the Plex client the server uses to
// connect to and execute queries against Plex
func initPlex(c *config.Config) {
	if plexClient != nil {
		return
	}
	plexClient = plex.New(c.Plex.Token, c.Plex.Address)
}

// initLLM creates the Ollama LLM client the server uses
// to connect to and execute generation
func initLLM(c *config.Config) error {
	if ollamaLlm != nil {
		return nil
	}
	var err error
	ollamaLlm, err = langchain.InitOllama(c.Ollama.Address, c.Ollama.Model)
	if err != nil {
		return err
	}
	return nil
}
