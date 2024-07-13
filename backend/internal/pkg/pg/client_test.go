package pg

import (
	"encoding/base64"
	"testing"
)

func TestQueryOptions(t *testing.T) {
	// Written mostly with Gemini
	var (
		testTitles    = []string{"test", "title"}
		anotherTitle  = []string{"another", "title"}
		stringTitles  = buildStringFromSlice(testTitles)
		anotherString = buildStringFromSlice(anotherTitle)
	)
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
			options:  []QueryOption{WithInputTitles(testTitles)},
			expected: queryOption{response: "", input: base64.StdEncoding.EncodeToString([]byte(stringTitles))},
		},
		{
			name:     "With Response",
			options:  []QueryOption{WithResponse("data")},
			expected: queryOption{response: "data", input: ""},
		},
		{
			name:     "With Both Options",
			options:  []QueryOption{WithInputTitles(anotherTitle), WithResponse("result")},
			expected: queryOption{response: "result", input: base64.StdEncoding.EncodeToString([]byte(anotherString))},
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
