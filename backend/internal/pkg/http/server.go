package httpinternal

import (
	"context"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"log"
	"net/http"

	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/config"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/langchain"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/pg"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/weaviate"
)

var (
	plexClient     *plex.PlexClient
	ollamaLlm      *ollama.LLM
	ollamaEmbedder *ollama.LLM
)

// StartServer initializes dependent services that are
// required to handle HTTP requests. This is blocking.
func StartServer(ctx context.Context, c *config.Config, shutdownChan chan error) {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Start Server"), telemetry.WithSpanPackage("httpinternal"))
	defer span.End()
	initPlex(ctx, c)
	if err := initLLM(ctx, c); err != nil {
		panic("could not initialize llms: " + err.Error())
	}
	if err := initVectorStore(ctx); err != nil {
		panic("could not init vector store: " + err.Error())
	}
	if err := initCacheStore(ctx); err != nil {
		panic("could not init cache store: " + err.Error())
	}
	initHttpServer(shutdownChan)
}

// initHttpServer is a blocking function that runs the HTTP server
// and passes requests to their respective handler functions
func initHttpServer(s chan error) {
	mux := http.NewServeMux()

	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	// Register handlers.
	handleFunc(recommendationPathway, recommendationHandler)

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	log.Println("serving http...")
	if err := http.ListenAndServe(":8090", handler); err != nil {
		s <- err
	}
	s <- nil
}

// initPlex creates the Plex client the server uses to
// connect to and execute queries against Plex
func initPlex(ctx context.Context, c *config.Config) {
	log.Println("initialzing plex client...")
	defer log.Println("initialized")
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("Init Plex"))
	defer span.End()
	if plexClient != nil {
		span.SetStatus(codes.Ok, "Plex client initialized previously")
		return
	}
	plexClient = plex.New(c.Plex.Token, c.Plex.Address)
	span.SetStatus(codes.Ok, "Plex client initialized")
}

// initLLM creates the Ollama LLM client the server uses
// to connect to and execute generation and embeddings
func initLLM(ctx context.Context, c *config.Config) error {
	if ollamaLlm != nil {
		return nil
	}
	var err error
	ollamaLlm, ollamaEmbedder, err = langchain.InitOllama(ctx, c.Ollama.Address, c.Ollama.LanguageModel, c.Ollama.EmbeddingModel)
	if err != nil {
		return err
	}
	return nil
}

// initVectorStore connects to Weaviate for storing
// Plex data and related embeddings and performs
// any migrations required for startup.
func initVectorStore(ctx context.Context) error {
	if err := weaviate.InitWeaviate(ctx, plexClient, ollamaEmbedder); err != nil {
		return err
	}
	return nil
}

// initCacheStore connects to a database used for
// storing responses from the LLM and the inputs
// used to generate them.
func initCacheStore(ctx context.Context) error {
	return pg.InitPostgres(ctx)
}
