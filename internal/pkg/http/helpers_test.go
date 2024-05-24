package httpinternal

import (
	"reflect"
	"testing"
)

func TestBuildStringFromSlice(t *testing.T) {
	testCases := []struct {
		name     string
		input    []any
		expected string
	}{
		{
			name:     "Empty slice",
			input:    []any{},
			expected: "[]",
		},
		{
			name:     "Slice of integers",
			input:    []any{1, 2, 3},
			expected: "[1 2 3]",
		},
		{
			name:     "Slice of strings",
			input:    []any{"hello", "world"},
			expected: `[hello world]`,
		},
		{
			name:     "Slice of mixed types",
			input:    []any{1, "foo", 3.14, true},
			expected: `[1 foo 3.14 true]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := buildStringFromSlice(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("For input %v, expected '%s', but got '%s'", tc.input, tc.expected, result)
			}
		})
	}
}
