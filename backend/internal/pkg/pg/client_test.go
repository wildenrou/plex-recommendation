package pg

import (
	"encoding/base64"
	"testing"
)

func TestQueryOptions(t *testing.T) {
	// Written entirely with Gemini
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
			name:     "With Input Title",
			options:  []QueryOption{WithInputTitles("test title")},
			expected: queryOption{response: "", input: base64.StdEncoding.EncodeToString([]byte("test title"))},
		},
		{
			name:     "With Response",
			options:  []QueryOption{WithReponse("data")},
			expected: queryOption{response: "data", input: ""},
		},
		{
			name:     "With Both Options",
			options:  []QueryOption{WithInputTitles("another title"), WithReponse("result")},
			expected: queryOption{response: "result", input: base64.StdEncoding.EncodeToString([]byte("another title"))},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			qo := queryOption{}
			for _, opt := range tc.options {
				opt(&qo)
			}

			if qo.input != tc.expected.input || qo.response != tc.expected.response {
				t.Errorf("Expected: %v, Got: %v", tc.expected, qo)
			}
		})
	}
}
