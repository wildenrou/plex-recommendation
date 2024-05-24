package plex

import (
	"log"
	"net/http"

	httpinternal "github.com/wgeorgecook/plex-recommendation/internal/pkg/http"
)

type Client interface {
	Connect(string, bool) string
	MakeNetworkRequest(string, string) (*http.Response, error)
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
		httpClient:  httpinternal.NewClient(),
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
func (pc PlexClient) MakeNetworkRequest(endpoint, method string) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
