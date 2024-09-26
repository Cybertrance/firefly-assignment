package article

import (
	"firefly-assignment/config"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// extractArticleContent parses the HTML document to extract the article content
func extractArticleContent(body string) (string, error) {
	// Create a goquery document from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error loading HTML: %v", err)
	}

	// Find the article content. For Engadget, the article content is inside <div> with class `caas-body`
	var containerSelector = config.AppConfig.ContainerSelector
	articleContent := doc.Find(containerSelector)
	if articleContent.Length() == 0 {
		return "", fmt.Errorf("could not find article content")
	}

	// Extract and return the text content
	return articleContent.Text(), nil
}

func getWords(article string) []string {
	return strings.Fields(article)
}

func GetArticleWords(rawBody string) ([]string, error) {
	article, err := extractArticleContent(rawBody)

	if err != nil {
		return nil, err
	}

	return getWords(article), nil
}
