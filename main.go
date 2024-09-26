package main

import (
	"bufio"
	"context"
	"firefly-assignment/article"
	"firefly-assignment/config"
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

// Configuration settings
var (
	nResults              int
	sourceUrlFileName     string
	requestsPerSecond     rate.Limit
	burstSize             int
	maxConcurrentRequests int8
)

var (
	wg               sync.WaitGroup
	wordBankChannel  chan utils.WordBank    = make(chan utils.WordBank, 1)
	validWords       utils.WordBank         = make(utils.WordBank)
	wordFrequencyMap utils.WordFrequencyMap = make(utils.WordFrequencyMap)
	processedURLs    int32                  = 0
	erroredURLs      int32                  = 0

	semaphoreMaxConcRequests chan struct{}
)

// processURL processes a URL by first fetching the raw content from the URL, scraping an article and finally updating the word frequency map.
func processURL(url string) {
	// Use a semaphore (with size `maxConcRequests`) to limit the number of concurrent URLs processed.
	semaphoreMaxConcRequests <- struct{}{}
	defer func() { <-semaphoreMaxConcRequests }()

	defer wg.Done()

	fmt.Printf("\n[INFO] - Processing URL: %v", url)

	body, err := network.FetchContent(url)
	if err != nil {
		fmt.Printf("\n[ERROR] - Failed to fetch URL: %v with Error: %v", url, err)
		erroredURLs++
		return
	}

	articleWords, err := article.GetArticleWords(body)
	if err != nil {
		fmt.Printf("\n[ERROR] - Failed to extract article content for URL: %v with Error: %v", url, err)
		erroredURLs++
		return
	}

	// Initialize the word bank of valid words.
	if len(validWords) == 0 {
		validWords = <-wordBankChannel
	}

	wordOps.CountWords(articleWords, validWords, wordFrequencyMap)
	processedURLs++
}

// getURLsFromFile gets the URLs for the articles to be scraped from the configured file
func getURLsFromFile() ([]string, error) {
	file, err := os.Open("static/" + sourceUrlFileName)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	// Create a scanner to read the file line by line
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func main() {
	// Load config from 'config.yaml' if available.
	config.LoadConfig()

	// Set Configuration settings
	nResults = config.AppConfig.TopResults
	sourceUrlFileName = config.AppConfig.SourceURLFileName
	requestsPerSecond = config.AppConfig.RequestsPerSecond
	burstSize = config.AppConfig.BurstSize
	maxConcurrentRequests = config.AppConfig.MaxConcurrentRequests
	semaphoreMaxConcRequests = make(chan struct{}, maxConcurrentRequests)

	// 1. Initialize the word bank of valid words
	go wordBank.Initialize(wordBankChannel)

	// 2. Get the URLs from file
	urls, err := getURLsFromFile()

	if err != nil {
		log.Fatalln("[ERROR] - No URLs to fetch content from")
	}

	limiter := rate.NewLimiter(requestsPerSecond, burstSize)

	// 3. For each URL, scrape and process the data.
	for _, url := range urls {
		limiter.Wait(context.Background())
		wg.Add(1)
		go processURL(url)
	}
	wg.Wait()

	// 4. Get Top N words
	var topNWords = wordOps.GetTopNWords(nResults, wordFrequencyMap)

	fmt.Printf("\n\n========")
	fmt.Printf("\nTotal entries: %v", len(urls))
	fmt.Printf("\nProcessed entries: %v", processedURLs)
	fmt.Printf("\nErrored entries: %v", erroredURLs)
	fmt.Println("\nTop 10 words:")
	fmt.Println(display.GetPrettyJSON(topNWords))
}
