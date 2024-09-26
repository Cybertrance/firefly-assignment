package wordOps

import (
	"firefly-assignment/utils"
	"sort"
	"strings"
	"sync"
)

var mutex sync.Mutex

// GetTopWords sorts and extracts the top 'n' words from the word frequency map
func GetTopWords(n int, wordFrequencyMap utils.WordFrequencyMap) []utils.WordFreq {
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

// CountWords processes an article by splitting the article into its constituent words and updating the wordFrequencyMap if it exists in the wordBankMap.
func CountWords(articleWords []string, wordBank utils.WordBank, wordFrequencyMap utils.WordFrequencyMap) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, word := range articleWords {
		normalizedWord := strings.ToLower(word)
		if _, exists := wordBank[normalizedWord]; exists {
			wordFrequencyMap[normalizedWord]++
		}
	}
}
