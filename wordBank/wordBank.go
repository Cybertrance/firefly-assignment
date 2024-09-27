// Package wordBank provides functionality for managing a word bank,
// including fetching and validating words from an external source.
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

// Initialize fetches the word bank from a configured URL, filters the words based on
// predefined rules (e.g., words longer than 3 characters and composed of letters),
// and sends the result through the provided channel.
//
// Parameters:
//   - wordBankChannel: A channel to which the validated word bank will be sent.
//
// Returns:
//   - error: An error if the word bank cannot be fetched or processed.
func Initialize(wordBankChannel chan utils.WordBank) error {
	wordBankMap := make(utils.WordBank)
	wordBankURL = config.AppConfig.WordBankURL
	resp, err := http.Get(wordBankURL)
	if err != nil {
		log.Fatalf("[ERROR] - Error fetching wordbank source %v - %v", wordBankURL, err)
	}
	defer resp.Body.Close()

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
	return nil
}
