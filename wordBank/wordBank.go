package wordBank

import (
	"firefly-assignment/config"
	"firefly-assignment/utils"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

var (
	wordBankURL string
)

// Initialize initializes the internal word bank by fetching it from the source URL and validating based on some rules.
func Initialize(wordBankChannel chan utils.WordBank) error {
	wordBankMap := make(utils.WordBank)
	wordBankURL = config.AppConfig.WordBankURL
	resp, err := http.Get(wordBankURL)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("[ERROR] - Error fetching wordbank source %v - %v", word_bank_url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[ERROR] - Error reading response, err")
	}
	words := strings.Fields(string(body))

	for _, word := range words {
		if utf8.RuneCountInString(word) > 3 && utils.IsLetter(word) {
			wordBankMap[strings.ToLower(word)] = struct{}{}
		}
	}

	wordBankChannel <- wordBankMap
}
