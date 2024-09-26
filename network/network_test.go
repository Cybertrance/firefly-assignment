package network

import (
	"fmt"
	"testing"

	"github.com/valyala/fasthttp"
)

// Mock fasthttp.Client to simulate different responses
type mockClient struct {
	statusCode int
	body       string
	redirect   string
	err        error
}

func (m *mockClient) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	if m.err != nil {
		return m.err
	}
	resp.SetStatusCode(m.statusCode)
	if m.body != "" {
		resp.SetBody([]byte(m.body))
	}
	if m.redirect != "" {
		resp.Header.Set("Location", m.redirect)
	}
	return nil
}

func TestFetchContent(t *testing.T) {
	tests := []struct {
		name          string
		mockClient    *mockClient
		url           string
		expectedBody  string
		expectedError bool
	}{
		{
			name: "Valid content fetch",
			mockClient: &mockClient{
				statusCode: fasthttp.StatusOK,
				body:       "This is the body content",
			},
			url:           "http://example.com",
			expectedBody:  "This is the body content",
			expectedError: false,
		},
		{
			name: "Too many redirects",
			mockClient: &mockClient{
				statusCode: fasthttp.StatusMovedPermanently,
				redirect:   "http://redirect.com",
			},
			url:           "http://example.com",
			expectedBody:  "",
			expectedError: true,
		},
		{
			name: "Server blocking with status code 999",
			mockClient: &mockClient{
				statusCode: 999,
			},
			url:           "http://blocked.com",
			expectedBody:  "",
			expectedError: true,
		},
		{
			name: "Retry limit reached",
			mockClient: &mockClient{
				statusCode: fasthttp.StatusNotFound,
			},
			url:           "http://retry.com",
			expectedBody:  "",
			expectedError: true,
		},
		{
			name: "Network error",
			mockClient: &mockClient{
				err: fmt.Errorf("network error"),
			},
			url:           "http://network-error.com",
			expectedBody:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Replace the default fasthttp client with the mock client
			httpClient = tt.mockClient

			// Call the function being tested
			body, err := FetchContent(tt.url)

			// Check if the error matches the expectation
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			// Check if the returned body matches the expected body
			if body != tt.expectedBody {
				t.Errorf("expected body: %v, got: %v", tt.expectedBody, body)
			}
		})
	}
}
