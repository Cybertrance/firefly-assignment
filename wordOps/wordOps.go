package wordOps

import (
	"container/heap"
	"firefly-assignment/minheap"
	"firefly-assignment/utils"
	"strings"
	"sync"
)

var mutex sync.Mutex

// GetTopNWords sorts and extracts the top 'n' words from the word frequency map
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

	// Get top N words
	// var topWords []utils.WordFreq
	// for i := 0; i < n && i < len(wordList); i++ {
	// 	topWords = append(topWords, utils.WordFreq{Word: wordList[i].Word, Frequency: wordList[i].Frequency})
	// }
	// return topWords
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
