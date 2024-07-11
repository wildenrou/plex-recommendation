package plex

import (
	"context"
	"github.com/wgeorgecook/plex-recommendation/internal/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"log"
	"net/http"
	"time"
)

type Client interface {
	Connect(string, bool) string
	MakeNetworkRequest(context.Context, string, string) (*http.Response, error)
}

type PlexClient struct {
	accessToken string
	address     string
	httpClient  *http.Client
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

// Connect returns a string with our address and token
// for the provided sectionId
func (pc PlexClient) Connect(sectionID string, allMovies bool) string {
	log.Println("generating connection string")
	endpoint := "/all"
	if !allMovies {
		endpoint = "/recentlyViewed"
	}
	return "http://" +
		pc.address +
		":32400/library/sections/" +
		sectionID +
		endpoint +
		"?X-Plex-Token=" +
		pc.accessToken
}

// MakeNetworkRequest makes an HTTP request with the provided method
// to the provided endpoint
func (pc PlexClient) MakeNetworkRequest(ctx context.Context, endpoint, method string) (*http.Response, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "MakeNetworkRequest")
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
