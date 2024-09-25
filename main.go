package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

// Consts
// TODO: Export these to a configuration file later.
const (
	maxRedirects      = 10
	filePath          = "static/endg-urls"
	containerSelector = ".caas-body"
)

var (
	wg       sync.WaitGroup
	mutex    sync.Mutex
	articles []string
)

func getURLsFromFile() ([]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // Append each line to the slice
	}

	// Check for any errors encountered while reading the file
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// fetchContent makes a GET request to the given URL and returns the body content as a string
func fetchContent(url string) (string, error) {
	client := &fasthttp.Client{
		ReadBufferSize: 8192, // Increase buffer size to handle larger request header sizes.
	}

	redirectCount := 0

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
		err := client.Do(req, resp)
		if err != nil {
			return "", fmt.Errorf("error fetching URL: %v", err)
		}

		// Check if it's a redirect status code (301, 302, 303, 307, 308)
		statusCode := resp.StatusCode()
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

// extractArticleContent parses the HTML document to extract the article content
func extractArticleContent(body string) (string, error) {
	// Create a goquery document from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error loading HTML: %v", err)
	}

	// Find the article content. For Engadget, the article content is inside <div> with class `caas-body`
	articleContent := doc.Find(containerSelector)
	if articleContent.Length() == 0 {
		return "", fmt.Errorf("could not find article content")
	}

	// Extract and return the text content
	return articleContent.Text(), nil
}

func processArticle(url string) {
	defer wg.Done()

	body, err := fetchContent(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}

	article, err := extractArticleContent(body)
	if err != nil {
		log.Fatalf("Failed to extract article content: %v", err)
	}

	mutex.Lock()
	articles = append(articles, article)
	mutex.Unlock()
}

func main() {
	urls, error := getURLsFromFile()
	if error != nil {
		fmt.Println(error)
	}
	for _, url := range urls {
		wg.Add(1)

		go processArticle(url)

		if error != nil {
			fmt.Printf("Could not process article: %v", url)
		}
	}
	wg.Wait()
	fmt.Println(len(articles))
}
