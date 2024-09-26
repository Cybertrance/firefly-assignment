package wordBank

import (
	"firefly-assignment/config"
	"firefly-assignment/utils"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

var (
	wordBankURL string
)

// Initialize initializes the internal word bank by fetching it from the source URL and validating based on some rules.
func Initialize(wordBankChannel chan utils.WordBank) {
	wordBankMap := make(utils.WordBank)
	wordBankURL = config.AppConfig.WordBankURL
	resp, err := http.Get(wordBankURL)
	if err != nil {
		fmt.Println("Error fetching wordbank source:")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	words := strings.Fields(string(body))

	for _, word := range words {
		if utf8.RuneCountInString(word) > 3 && utils.IsLetter(word) {
			wordBankMap[strings.ToLower(word)] = struct{}{}
		}
	}

	wordBankChannel <- wordBankMap
}
