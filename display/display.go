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
		fmt.Println("Could not convert to JSON")
	}

	return string(prettyJSON), nil

}
