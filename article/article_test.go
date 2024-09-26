package article

import (
	"testing"
)

func TestGetArticleWords(t *testing.T) {
	tests := []struct {
		name          string
		inputHTML     string
		expectedWords []string
		expectedError bool
	}{
		{
			name:          "Valid article content",
			inputHTML:     `<html><body><div class="caas-body">This is a test article with some words.</div></body></html>`,
			expectedWords: []string{"This", "is", "a", "test", "article", "with", "some", "words."},
			expectedError: false,
		},
		{
			name:          "Empty article content",
			inputHTML:     `<html><body><div class="caas-body"></div></body></html>`,
			expectedWords: []string{},
			expectedError: false,
		},
		{
			name:          "Missing article content",
			inputHTML:     `<html><body><div class="wrong-class">No article content here.</div></body></html>`,
			expectedWords: nil,
			expectedError: true,
		},
		{
			name:          "Malformed HTML",
			inputHTML:     `<html><body><div class="caas-body>This is broken HTML</div></body></html>`,
			expectedWords: nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words, err := GetArticleWords(tt.inputHTML)

			// Check if the error matches the expectation
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			// If no error, check if the returned words match the expected result
			if err == nil && !equal(words, tt.expectedWords) {
				t.Errorf("expected words: %v, got: %v", tt.expectedWords, words)
			}
		})
	}
}

// Helper function to compare two slices of strings
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
