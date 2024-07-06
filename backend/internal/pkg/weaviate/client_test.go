package weaviate

import (
	"testing"

	"github.com/wgeorgecook/plex-recommendation/internal/pkg/plex"
)

// Written entirely by Gemini
func TestQueryOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []QueryOption
		expected queryOption
	}{
		{
			name:     "Empty Options",
			options:  []QueryOption{},
			expected: queryOption{},
		},
		{
			name:     "With Class Name",
			options:  []QueryOption{WithClassName("document")},
			expected: queryOption{className: "document", limit: 0},
		},
		{
			name:     "With Limit",
			options:  []QueryOption{WithLimit(10)},
			expected: queryOption{className: "", limit: 10},
		},
		{
			name:     "With Both Options",
			options:  []QueryOption{WithClassName("image"), WithLimit(25)},
			expected: queryOption{className: "image", limit: 25},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			qo := queryOption{}
			for _, opt := range tc.options {
				opt(&qo)
			}

			if qo.className != tc.expected.className || qo.limit != tc.expected.limit {
				t.Errorf("Expected: %v, Got: %v", tc.expected, qo)
			}
		})
	}
}

// Written with scaffolded code from Gemini
func TestWithVideos(t *testing.T) {
	const (
		kiki   = "Kiki's Delivery Service"
		totoro = "My Neighbor Totoro"
	)

	tests := []struct {
		name     string
		videos   []plex.VideoShort
		expected insertOption
	}{
		{
			name:     "Empty Videos",
			videos:   []plex.VideoShort{},
			expected: insertOption{videos: []plex.VideoShort{}},
		},
		{
			name:     "With Single Video",
			videos:   []plex.VideoShort{{Title: kiki}},
			expected: insertOption{videos: []plex.VideoShort{{Title: kiki}}},
		},
		{
			name:     "With Multiple Videos",
			videos:   []plex.VideoShort{{Title: kiki}, {Title: totoro}},
			expected: insertOption{videos: []plex.VideoShort{{Title: kiki}, {Title: totoro}}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			io := insertOption{}
			WithVideos(tc.videos)(&io)

			if len(io.videos) != len(tc.expected.videos) {
				t.Errorf("Expected video count: %d, Got: %d", len(tc.expected.videos), len(io.videos))
				return
			}

			for i, v := range io.videos {
				if v.Title != tc.expected.videos[i].Title {
					t.Errorf("Expected video title at index %d: %s, Got: %s", i, tc.expected.videos[i].Title, v.Title)
				}
			}
		})
	}
}
