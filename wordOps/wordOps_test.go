package wordOps

import (
	"firefly-assignment/utils"
	"testing"
)

func TestGetTopWords(t *testing.T) {
	tests := []struct {
		name          string
		n             int
		wordFrequency utils.WordFrequencyMap
		expected      []utils.WordFreq
	}{
		{
			name: "Get top 2 words",
			n:    2,
			wordFrequency: utils.WordFrequencyMap{
				"apple":  5,
				"banana": 3,
				"cherry": 1,
			},
			expected: []utils.WordFreq{
				{Word: "apple", Frequency: 5},
				{Word: "banana", Frequency: 3},
			},
		},
		{
			name: "Get top 1 word",
			n:    1,
			wordFrequency: utils.WordFrequencyMap{
				"apple":  2,
				"banana": 4,
			},
			expected: []utils.WordFreq{
				{Word: "banana", Frequency: 4},
			},
		},
		{
			name: "Get more words than available",
			n:    5,
			wordFrequency: utils.WordFrequencyMap{
				"apple":  2,
				"banana": 3,
			},
			expected: []utils.WordFreq{
				{Word: "banana", Frequency: 3},
				{Word: "apple", Frequency: 2},
			},
		},
		{
			name:          "Empty frequency map",
			n:             2,
			wordFrequency: utils.WordFrequencyMap{},
			expected:      []utils.WordFreq{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTopNWords(tt.n, tt.wordFrequency)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d words, got %d", len(tt.expected), len(result))
			}
			for i, wordFreq := range result {
				if wordFreq.Word != tt.expected[i].Word || wordFreq.Frequency != tt.expected[i].Frequency {
					t.Errorf("expected %v, got %v", tt.expected[i], wordFreq)
				}
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name            string
		articleWords    []string
		wordBank        utils.WordBank
		initialFreqMap  utils.WordFrequencyMap
		expectedFreqMap utils.WordFrequencyMap
	}{
		{
			name:         "Count words present in word bank",
			articleWords: []string{"Apple", "Banana", "Apple"},
			wordBank: utils.WordBank{
				"apple":  struct{}{},
				"banana": struct{}{},
			},
			initialFreqMap:  utils.WordFrequencyMap{},
			expectedFreqMap: utils.WordFrequencyMap{"apple": 2, "banana": 1},
		},
		{
			name:         "Count words with non-existing words in word bank",
			articleWords: []string{"Apple", "Orange", "Cherry"},
			wordBank: utils.WordBank{
				"apple": struct{}{},
			},
			initialFreqMap:  utils.WordFrequencyMap{},
			expectedFreqMap: utils.WordFrequencyMap{"apple": 1},
		},
		{
			name:            "Empty word bank",
			articleWords:    []string{"Apple", "Banana"},
			wordBank:        utils.WordBank{},
			initialFreqMap:  utils.WordFrequencyMap{},
			expectedFreqMap: utils.WordFrequencyMap{},
		},
		{
			name:         "Word frequency map with initial values",
			articleWords: []string{"Banana", "Banana"},
			wordBank: utils.WordBank{
				"banana": struct{}{},
			},
			initialFreqMap:  utils.WordFrequencyMap{"banana": 1},
			expectedFreqMap: utils.WordFrequencyMap{"banana": 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CountWords(tt.articleWords, tt.wordBank, tt.initialFreqMap)
			if len(tt.initialFreqMap) != len(tt.expectedFreqMap) {
				t.Fatalf("expected %d entries in frequency map, got %d", len(tt.expectedFreqMap), len(tt.initialFreqMap))
			}
			for word, freq := range tt.expectedFreqMap {
				if tt.initialFreqMap[word] != freq {
					t.Errorf("expected frequency of word %q to be %d, got %d", word, freq, tt.initialFreqMap[word])
				}
			}
		})
	}
}
