/*
Package display provides utilities for formatting and displaying data in various formats,
such as pretty-formatted JSON. It is designed to help convert complex data structures
into human-readable formats, making it useful for logging, debugging, or presenting data
in a clear and structured manner.
*/
package display

import (
	"encoding/json"
	"firefly-assignment/utils"
	"fmt"
)

// GetPrettyJSON takes a slice of WordFreq structs and returns a pretty-formatted JSON string.
// The JSON is indented for readability with four spaces per indentation level.
//
// Parameters:
//   - words: A slice of utils.WordFreq structs representing word frequencies.
//
// Returns:
//   - string: A pretty-formatted JSON string representing the word frequencies.
//   - error: An error if the struct cannot be converted to JSON.
func GetPrettyJSON(words []utils.WordFreq) (string, error) {
	prettyJSON, err := json.MarshalIndent(words, "", "    ")

	if err != nil {
		return "", fmt.Errorf("[ERROR] - Could not convert struct to JSON - %w", err)
	}

	return string(prettyJSON), nil

}
