package display

import (
	"encoding/json"
	"firefly-assignment/utils"
	"fmt"
)

// GetPrettyJSON gets a pretty-json from the words struct
func GetPrettyJSON(words []utils.WordFreq) (string, error) {
	prettyJSON, err := json.MarshalIndent(words, "", "    ")

	if err != nil {
		return "", fmt.Errorf("[ERROR] - Could not convert struct to JSON - %w", err)
	}

	return string(prettyJSON), nil

}
