package main

import (
	"bufio"
	"context"
	"firefly-assignment/article"
	"firefly-assignment/display"
	"firefly-assignment/network"
	"firefly-assignment/utils"
	"firefly-assignment/wordBank"
	"firefly-assignment/wordOps"
	"fmt"
	"log"
	"os"
	"sync"

	"golang.org/x/time/rate"
)

// Consts
// TODO: Export these to a configuration file later.
const (
	topResults            = 10
	filePath              = "static/endg-urls"
	requestsPerSecond     = 20
	burstSize             = 20
	maxConcurrentRequests = 20
)

var (
	wg               sync.WaitGroup
	wordBankChannel                         = make(chan utils.WordBank, 1)
	validWords       utils.WordBank         = make(utils.WordBank)
	wordFrequencyMap utils.WordFrequencyMap = make(utils.WordFrequencyMap)
	processedURLs    int32                  = 0
	erroredURLs      int32                  = 0

	semaphoreMaxConcRequests = make(chan struct{}, maxConcurrentRequests)
)

// processURL processes a URL by first fetching the raw content, scraping the article text and finally updating the word frequency map.
// Use a semaphore (with size `maxConcRequests`) to limit the number of concurrent URLs processed.
func processURL(url string) {
	semaphoreMaxConcRequests <- struct{}{}
	defer func() { <-semaphoreMaxConcRequests }()

	defer wg.Done()

	fmt.Printf("\nProcessing URL: %v", url)

	body, err := network.FetchContent(url)
	if err != nil {
		fmt.Printf("\nFailed to fetch URL: %v with Error: %v", url, err)
		erroredURLs++
		return
	}

	articleWords, err := article.GetArticleWords(body)
	if err != nil {
		fmt.Printf("\nFailed to extract article content for URL: %v with Error: %v", url, err)
		erroredURLs++
		return
	}

	if len(validWords) == 0 {
		validWords = <-wordBankChannel
	}

	wordOps.CountWords(articleWords, validWords, wordFrequencyMap)
	processedURLs++
}

// getURLsFromFile gets the URLs for the articles to be scraped from the 'endg-urls' file
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

func main() {
	go wordBank.Initialize(wordBankChannel)

	urls, error := getURLsFromFile()

	if error != nil {
		log.Fatalln("No URLs to fetch content from. ERROR: ", error)
	}

	limiter := rate.NewLimiter(requestsPerSecond, burstSize)

	for _, url := range urls {
		limiter.Wait(context.Background())
		wg.Add(1)
		go processURL(url)
	}
	wg.Wait()

	var topNWords = wordOps.GetTopWords(topResults, wordFrequencyMap)

	fmt.Printf("\n\n========")
	fmt.Printf("\nTotal entries: %v", len(urls))
	fmt.Printf("\nProcessed entries: %v", processedURLs)
	fmt.Printf("\nErrored entries: %v", erroredURLs)
	fmt.Println("\nTop 10 words:")
	fmt.Println(display.GetPrettyJSON(topNWords))
}
