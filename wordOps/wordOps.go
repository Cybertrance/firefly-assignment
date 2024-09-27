// Package wordOps provides operations for processing word frequencies,
// including counting word occurrences and extracting the top N frequent words.
package wordOps

import (
	"container/heap"
	"firefly-assignment/minheap"
	"firefly-assignment/utils"
	"strings"
	"sync"
)

var mutex sync.Mutex

// GetTopNWords returns the top 'n' words with the highest frequencies from the given word frequency map.
// It uses a min-heap to efficiently keep track of the top words.
//
// Parameters:
//   - n: The number of top words to return.
//   - wordFrequencyMap: A map where keys are words and values are their frequencies.
//
// Returns:
//   - []utils.WordFreq: A slice containing the top 'n' words with their frequencies, sorted by frequency.
func GetTopNWords(n int, wordFrequencyMap utils.WordFrequencyMap) []utils.WordFreq {
	mutex.Lock()
	defer mutex.Unlock()

	// Initialize heatmap
	h := &minheap.MinHeap{}
	heap.Init(h)

	for word, count := range wordFrequencyMap {
		if h.Len() < n {
			heap.Push(h, utils.WordFreq{Word: word, Frequency: count})
		} else if count > (*h)[0].Frequency {
			heap.Pop(h)
			heap.Push(h, utils.WordFreq{Word: word, Frequency: count})
		}
	}

	// Get the items from the heap.
	result := make([]utils.WordFreq, 0, n)
	for h.Len() > 0 {
		result = append(result, heap.Pop(h).(utils.WordFreq))
	}

	// Reverse result to show highest frequency first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// CountWords updates the word frequency map by counting occurrences of words in the article that exist in the word bank.
//
// Parameters:
//   - articleWords: A slice of words from the article to be processed.
//   - wordBank: A set of valid words used for filtering the article words.
//   - wordFrequencyMap: A map where word counts will be updated.
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
