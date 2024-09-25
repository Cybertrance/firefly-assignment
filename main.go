package main

import (
	"bufio"
	"firefly-assignment/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

// Consts
// TODO: Export these to a configuration file later.
const (
	maxRedirects      = 10
	nTop              = 10
	filePath          = "static/endg-urls"
	containerSelector = ".caas-body"
	wordBankURL       = "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt"
)

type WordFreq struct {
	word      string
	frequency int32
}

var (
	wg               sync.WaitGroup
	mutex            sync.Mutex
	wordFrequencyMap map[string]int32    = make(map[string]int32)
	wordBankMap      map[string]struct{} = make(map[string]struct{})
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

// countWords processes an article by splitting the article into its constituent words and updating the wordFrequencyMap if it exists in the wordBankMap.
func countWords(article string) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, word := range strings.Fields(article) {
		normalizedWord := strings.ToLower(word)
		if _, exists := wordBankMap[normalizedWord]; exists {
			wordFrequencyMap[normalizedWord]++
		}
	}
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

	countWords(article)
}

// initializeWordBank initializes the internal word bank by fetching it from the source URL and validating based on some rules.
func initializeWordBank() {
	resp, err := http.Get(wordBankURL)
	if err != nil {
		fmt.Println("Error fetching wordbank source:")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	words := strings.Fields(string(body))

	for _, word := range words {
		if utf8.RuneCountInString(word) > 3 && utils.IsLetter(word) {
			wordBankMap[strings.ToLower(word)] = struct{}{}
		}
	}

	fmt.Println(len(wordBankMap))
}

// getTopWords sorts and extracts the top 'n' words from the word frequency map
func getTopWords(n int) []WordFreq {
	mutex.Lock()
	defer mutex.Unlock()

	// Convert map to slice of pairs
	var wordList []WordFreq
	for word, count := range wordFrequencyMap {
		wordList = append(wordList, WordFreq{word, count})
	}

	// Sort by frequency
	sort.Slice(wordList, func(i, j int) bool {
		return wordList[i].frequency > wordList[j].frequency
	})

	// Get top N words
	var topWords []WordFreq
	for i := 0; i < n && i < len(wordList); i++ {
		topWords = append(topWords, WordFreq{wordList[i].word, wordList[i].frequency})
	}
	return topWords
}

func main() {
	initializeWordBank()

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
	for k, wordFreq := range wordFrequencyMap {
		fmt.Println(k, wordFreq)
	}

	topNWords := getTopWords(10)

	for word, freq := range topNWords {
		fmt.Println(word, freq)
	}
}
