package display

import (
	"encoding/json"
	"firefly-assignment/utils"
	"fmt"
)

func GetPrettyJSON(words []utils.WordFreq) (string, error) {
	prettyJSON, err := json.MarshalIndent(words, "", "    ")

	if err != nil {
		fmt.Println("Could not convert to JSON")
	}

	return string(prettyJSON), nil

}
