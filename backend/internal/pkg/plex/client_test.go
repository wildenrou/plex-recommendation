package plex

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// TestPlexClientConnect was written entirely
// by an LLM with minor input refinements.
func TestPlexClientConnect(t *testing.T) {
	testClient := PlexClient{
		address:     "localhost",
		accessToken: "random_token_value",
	}

	testCases := []struct {
		name        string
		sectionID   string
		allMovies   bool
		expectedURL string
	}{
		{
			name:        "All Movies",
			sectionID:   "123",
			allMovies:   true,
			expectedURL: "http://localhost:32400/library/sections/123/all?X-Plex-Token=random_token_value",
		},
		{
			name:        "Recently Viewed",
			sectionID:   "456", // Different section ID for variety
			allMovies:   false,
			expectedURL: "http://localhost:32400/library/sections/456/recentlyViewed?X-Plex-Token=random_token_value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualURL := testClient.Connect(tc.sectionID, tc.allMovies)
			if actualURL != tc.expectedURL {
				t.Errorf("URL mismatch for %s: expected %s, got %s", tc.name, tc.expectedURL, actualURL)
			}
		})
	}
}

// TestPlexClientMakeNetworkRequest was written entirely
// with an LLM.
func TestPlexClientMakeNetworkRequest(t *testing.T) {
	testClient := New("randomToken", "localhost")

	testCases := []struct {
		name       string
		endpoint   string
		method     string
		mockStatus int
		mockBody   string
		mockError  error
		wantErr    bool
	}{
		{
			name:       "Successful GET",
			endpoint:   "http://localhost/test",
			method:     http.MethodGet,
			mockStatus: http.StatusOK,
			mockBody:   "Test response",
			wantErr:    false,
		},
		{
			name:       "Failed POST",
			endpoint:   "http://localhost/post",
			method:     http.MethodPost,
			mockStatus: http.StatusInternalServerError,
			mockError:  errors.New("Server error"),
			wantErr:    true,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock response
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tc.mockError != nil {
				httpmock.RegisterResponder(tc.method, tc.endpoint, httpmock.NewErrorResponder(tc.mockError))
			} else {
				httpmock.RegisterResponder(tc.method, tc.endpoint, httpmock.NewStringResponder(tc.mockStatus, tc.mockBody))
			}

			// Make the request
			resp, err := testClient.MakeNetworkRequest(context.Background(), tc.endpoint, tc.method)

			// Assertions
			if (err != nil) != tc.wantErr {
				t.Errorf("Unexpected error: expected %v, got %v", tc.wantErr, err)
			}

			if !tc.wantErr && resp.StatusCode != tc.mockStatus {
				t.Errorf("Status code mismatch: expected %d, got %d", tc.mockStatus, resp.StatusCode)
			}

			if !tc.wantErr && resp.Body != nil {
				defer resp.Body.Close()
				// Add body content verification if needed
			}
		})
	}
}
