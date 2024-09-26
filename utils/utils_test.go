package utils

import (
	"testing"
)

func TestIsLetter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "All letters",
			input:    "abcABC",
			expected: true,
		},
		{
			name:     "Mixed letters and digits",
			input:    "abc123",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: true, // technically no non-letter characters, so true
		},
		{
			name:     "String with spaces",
			input:    "abc ABC",
			expected: false,
		},
		{
			name:     "Non-letter Unicode characters",
			input:    "abc-äöü",
			expected: false,
		},
		{
			name:     "All Unicode letters",
			input:    "äöü",
			expected: true,
		},
		{
			name:     "Special characters",
			input:    "abc$%^",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLetter(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
