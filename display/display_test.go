package display

import (
	"firefly-assignment/utils"
	"testing"
)

func TestGetPrettyJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         []utils.WordFreq
		expectedError bool
	}{
		{
			name: "Valid input",
			input: []utils.WordFreq{
				{Word: "test", Frequency: 10},
				{Word: "example", Frequency: 5},
			},
			expectedError: false,
		},
		{
			name:          "Empty input",
			input:         []utils.WordFreq{},
			expectedError: false,
		},
		{
			name:          "Nil input",
			input:         nil,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetPrettyJSON(tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			// Only check if the result is not empty in successful cases
			if !tt.expectedError && result == "" {
				t.Error("expected non-empty JSON output but got an empty string")
			}
		})
	}
}
