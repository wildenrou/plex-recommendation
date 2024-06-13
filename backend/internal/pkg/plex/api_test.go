package plex

import (
	"testing"
)

// TestFullToShort was written entirely with
// an LLM
func TestFullToShort(t *testing.T) {
	testCases := []struct {
		name     string
		videos   []Video
		limit    int
		expected []VideoShort
	}{
		{
			name:     "Empty Input",
			videos:   []Video{},
			limit:    3,
			expected: []VideoShort{},
		},
		{
			name: "Limit Smaller Than Input",
			videos: []Video{
				{Title: "Video 1", Summary: "Summary 1", ContentRating: "G"},
				{Title: "Video 2", Summary: "Summary 2", ContentRating: "PG"},
				{Title: "Video 3", Summary: "Summary 3", ContentRating: "PG-13"},
			},
			limit: 2,
			expected: []VideoShort{
				{Title: "Video 1", Summary: "Summary 1", ContentRating: "G"},
				{Title: "Video 2", Summary: "Summary 2", ContentRating: "PG"},
			},
		},
		{
			name: "Limit Larger Than Input",
			videos: []Video{
				{Title: "Movie A", Summary: "Action-packed", ContentRating: "R"},
			},
			limit: 5,
			expected: []VideoShort{
				{Title: "Movie A", Summary: "Action-packed", ContentRating: "R"},
			},
		},
		// Add more test cases here if you want to cover other scenarios
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fullToShort(tc.videos, tc.limit)
			if len(result) != len(tc.expected) {
				t.Fatalf("Length mismatch: expected %d, got %d", len(tc.expected), len(result))
			}
			for i, short := range result {
				if short != tc.expected[i] {
					t.Errorf("Mismatch at index %d:\nExpected: %+v\nGot:      %+v", i, tc.expected[i], short)
				}
			}
		})
	}
}
