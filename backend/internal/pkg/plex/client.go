package plex

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Client interface {
	Connect(...ConnectOption) string
	GetDefaultLibrarySection() string
	MakeNetworkRequest(context.Context, string, string) (*http.Response, error)
}

type PlexClient struct {
	accessToken           string
	address               string
	httpClient            *http.Client
	defaultLibrarySection string
}

func New(accesstoken, address string) *PlexClient {
	return &PlexClient{
		accessToken: accesstoken,
		address:     address,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type connectOptions struct {
	sectionId string
	allMovies bool
}
type ConnectOption func(*connectOptions)

func WithSectionID(sectionID string) ConnectOption {
	return func(o *connectOptions) {
		o.sectionId = sectionID
	}
}

func WithAllMovies(allMovies bool) ConnectOption {
	return func(o *connectOptions) {
		o.allMovies = allMovies
	}
}

// Connect returns a string with our address and token
// for the provided sectionId
func (pc PlexClient) Connect(opts ...ConnectOption) string {
	log.Println("generating connection string")
	var options = connectOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	endpoint := "/all"
	if !options.allMovies {
		endpoint = "/recentlyViewed"
	}

	sectionId := pc.defaultLibrarySection
	if options.sectionId != "" {
		sectionId = options.sectionId
	}

	return "http://" +
		pc.address +
		":32400/library/sections/" +
		sectionId +
		endpoint +
		"?X-Plex-Token=" +
		pc.accessToken
}

// GetDefaultLibrarySection returns the library section to
// default to if none is provided on a request with section
// options.
func (pc PlexClient) GetDefaultLibrarySection() string {
	return pc.defaultLibrarySection
}

// MakeNetworkRequest makes an HTTP request with the provided method
// to the provided endpoint
func (pc PlexClient) MakeNetworkRequest(ctx context.Context, endpoint, method string) (*http.Response, error) {
	ctx, span := telemetry.StartSpan(ctx, telemetry.WithSpanName("MakeNetworkRequest"))
	defer span.End()
	span.SetAttributes(attribute.String("endpoint", endpoint))
	span.SetAttributes(attribute.String("method", method))
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.SetStatus(codes.Ok, resp.Status)
	return resp, nil
}
