package main

import (
	"bufio"
	"encoding/json"
	"firefly-assignment/article"
	"firefly-assignment/network"
	"firefly-assignment/utils"
	"firefly-assignment/wordBank"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

// Consts
// TODO: Export these to a configuration file later.
const (
	topResults = 10
	filePath   = "static/endg-urls"
)

var (
	wg               sync.WaitGroup
	mutex            sync.Mutex
	wordBankChannel                   = make(chan utils.WordBank, 1)
	validWords       utils.WordBank   = make(utils.WordBank)
	wordFrequencyMap map[string]int32 = make(map[string]int32)
)

// processURL processes a URL by first fetching the raw content, scraping the article text and finally updating the word frequency map.
func processURL(url string) {
	defer wg.Done()

	body, err := network.FetchContent(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}

	articleWords, err := article.GetArticleWords(body)
	if err != nil {
		log.Fatalf("Failed to extract article content: %v", err)
	}

	if len(validWords) == 0 {
		validWords = <-wordBankChannel
	}

	countWords(articleWords, validWords)
}

// CountWords processes an article by splitting the article into its constituent words and updating the wordFrequencyMap if it exists in the wordBankMap.
func countWords(articleWords []string, wordBank utils.WordBank) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, word := range articleWords {
		normalizedWord := strings.ToLower(word)
		if _, exists := wordBank[normalizedWord]; exists {
			wordFrequencyMap[normalizedWord]++
		}
	}
}

// getTopWords sorts and extracts the top 'n' words from the word frequency map
func getTopWords(n int) []utils.WordFreq {
	mutex.Lock()
	defer mutex.Unlock()

	// Convert map to slice of pairs
	var wordList []utils.WordFreq
	for word, count := range wordFrequencyMap {
		wordList = append(wordList, utils.WordFreq{Word: word, Frequency: count})
	}

	// Sort by frequency
	sort.Slice(wordList, func(i, j int) bool {
		return wordList[i].Frequency > wordList[j].Frequency
	})

	// Get top N words
	var topWords []utils.WordFreq
	for i := 0; i < n && i < len(wordList); i++ {
		topWords = append(topWords, utils.WordFreq{Word: wordList[i].Word, Frequency: wordList[i].Frequency})
	}
	return topWords
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

func getPrettyJSON(words []utils.WordFreq) (string, error) {
	prettyJSON, err := json.MarshalIndent(words, "", "    ")

	if err != nil {
		fmt.Println("Could not convert to JSON")
	}

	return string(prettyJSON), nil

}

func main() {
	go wordBank.Initialize(wordBankChannel)

	urls, error := getURLsFromFile()

	if error != nil {
		log.Fatalln("No URLs to fetch content from. ERROR: ", error)
	}

	for _, url := range urls {
		wg.Add(1)
		go processURL(url)

		if error != nil {
			fmt.Printf("Could not process article: %v", url)
		}
	}
	wg.Wait()

	// for k, wordFreq := range wordFrequencyMap {
	// 	fmt.Println(k, wordFreq)
	// }

	var topNWords = getTopWords(topResults)

	fmt.Println(getPrettyJSON(topNWords))
}
