package network

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

const (
	maxRetries   int8 = 10
	maxRedirects int8 = 10
)

// Creating an abstraction over fasthttp.Client to allow for mocking in test.
// HTTPClient defines the interface for an HTTP client
type HTTPClient interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
}

// DefaultHTTPClient wraps the default fasthttp.Client
type DefaultHTTPClient struct {
	Client *fasthttp.Client
}

func (c *DefaultHTTPClient) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	return c.Client.Do(req, resp)
}

var httpClient HTTPClient = &DefaultHTTPClient{
	Client: &fasthttp.Client{
		ReadBufferSize: 8192, // Increase buffer size to handle larger request header sizes.
	},
}

// fetchContent makes a GET request to the given URL and returns the body content as a string
func FetchContent(url string) (string, error) {

	var redirectCount int8 = 0
	var retryCount int8 = 0

	// Handle the edge-case where the URL could be a redirect.
	// Attempt to discover the redirected URL and fetch the content from there.
	for {
		// Manually create a fasthttp request and response to use with the custom client
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		// Set the URL for the request
		req.SetRequestURI(url)

		// Make the GET request
		err := httpClient.Do(req, resp)
		// if err != nil {
		// 	return "", fmt.Errorf("error fetching URL: %v", err)
		// }

		statusCode := resp.StatusCode()

		if statusCode == 999 {
			return "", fmt.Errorf("blocked by the server with status code 999")
		}

		if statusCode == 404 || err != nil {
			if retryCount >= maxRetries {
				return "", fmt.Errorf("too many retries")
			}

			fmt.Println("Retrying URL:" + url)
			retryCount++
			continue
		}

		// Check if it's a redirect status code (301, 302, 303, 307, 308)
		if statusCode >= 300 && statusCode < 400 {
			if redirectCount >= maxRedirects {
				return "", fmt.Errorf("too many redirects")
			}

			// Get the "Location" header to find the new URL
			newURL := resp.Header.Peek("Location")
			if newURL == nil {
				return "", fmt.Errorf("redirect with no Location header")
			}

			// Update the URL to the new location and continue the loop
			url = string(newURL)
			redirectCount++
			continue
		}

		// Return the response body as a string if it's not a redirect
		if statusCode == fasthttp.StatusOK {
			return string(resp.Body()), nil
		}

		// If the status code is not OK or a redirect, return an error
		return "", fmt.Errorf("received non-200 response: %d", statusCode)
	}
}
