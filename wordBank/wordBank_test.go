package wordBank

import (
	"firefly-assignment/config"
	"firefly-assignment/utils"
	"io"
	"net/http"
	"strings"
	"testing"
)

// Mock the http.Get function to simulate network calls
var httpGetFunc = func(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("apple banana cherry dog elephant")), // Mocked response
	}, nil
}

func mockHTTPGet(url string) (*http.Response, error) {
	return httpGetFunc(url)
}

func TestInitialize(t *testing.T) {
	config.LoadConfig()

	tests := []struct {
		name          string
		mockResponse  string
		expectedWords []string
		expectError   bool
	}{
		{
			name:          "Valid word bank response",
			mockResponse:  "apple banana cherry dog elephant",
			expectedWords: []string{"apple", "banana", "cherry", "elephant"}, // "dog" is ignored (<= 3 letters)
			expectError:   false,
		},
		{
			name:          "Empty response",
			mockResponse:  "",
			expectedWords: []string{},
			expectError:   false,
		},
		{
			name:          "Mixed valid and invalid words",
			mockResponse:  "apple banana 12345 $%@!",
			expectedWords: []string{"apple", "banana"}, // "12345" and "$%@!" are ignored
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the http.Get call
			httpGetFunc = func(url string) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
				}, nil
			}

			// Create a channel to receive the word bank
			wordBankChannel := make(chan utils.WordBank, 1)

			// Call Initialize function
			go Initialize(wordBankChannel)

			// Get the word bank from the channel
			wordBank := <-wordBankChannel

			// Verify the words in the word bank
			for _, word := range tt.expectedWords {
				if _, exists := wordBank[word]; !exists {
					t.Errorf("expected word %s in word bank, but it was not found", word)
				}
			}
		})
	}
}
