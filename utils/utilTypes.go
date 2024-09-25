package utils

type WordBank = map[string]struct{}

type WordFrequencyMap = map[string]int32

type WordFreq struct {
	Word      string
	Frequency int32
}
