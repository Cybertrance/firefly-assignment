/*
Package article provides utilities for extracting and processing article content
from raw HTML. The package is designed to work with HTML bodies of web articles,
extracting meaningful content, and further splitting that content into words for
analysis or other purposes.
*/
package article

import (
	"firefly-assignment/config"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// extractArticleContent extracts and returns the textual content of an article
// from the provided HTML string. It uses a CSS selector to identify the container
// of the article's body, which is configurable via the application's configuration.
//
// Parameters:
//   - body: A string containing the HTML content from which the article text will be extracted.
//
// Returns:
//   - string: The extracted article text.
//   - error: An error if the article content cannot be found or the HTML cannot be parsed.
func extractArticleContent(body string) (string, error) {
	// Create a goquery document from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("[ERROR] - error loading HTML: %v", err)
	}

	// Find the article content. For Engadget, the article content is inside <div> with class `caas-body`
	// However, you can configure this selector in the config ('container_selector').
	containerSelector := config.AppConfig.ContainerSelector
	articleContent := doc.Find(containerSelector)
	if articleContent.Length() == 0 {
		return "", fmt.Errorf("[ERROR] - could not find article content")
	}

	// Extract and return the text content
	return articleContent.Text(), nil
}

// getWords splits the provided article text into individual words.
// It uses whitespace as the delimiter to separate the words.
//
// Parameters:
//   - article: A string representing the article text to be split into words.
//
// Returns:
//   - []string: A slice containing individual words from the article text.
func getWords(article string) []string {
	return strings.Fields(article)
}

// GetArticleWords extracts the text content of an article from the provided raw HTML body
// and splits the article into individual words.
//
// This function first extracts the article content using extractArticleContent, and then
// splits the content into words using getWords.
//
// Parameters:
//   - rawBody: A string containing the raw HTML body of the article.
//
// Returns:
//   - []string: A slice containing individual words from the extracted article content.
//   - error: An error if the article content cannot be extracted.
func GetArticleWords(rawBody string) ([]string, error) {
	article, err := extractArticleContent(rawBody)

	if err != nil {
		return nil, err
	}

	return getWords(article), nil
}
