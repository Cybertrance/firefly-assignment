package config

import (
	"testing"

	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

func TestLoadConfig(t *testing.T) {

	tests := []struct {
		name             string
		mockConfig       string
		expectedConfig   Config
		shouldUseDefault bool
	}{
		{
			name:       "No config file, use defaults",
			mockConfig: "",
			expectedConfig: Config{
				TopResults:            10,
				SourceURLFileName:     "endg-urls",
				WordBankURL:           "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt",
				ContainerSelector:     ".caas-body",
				RequestsPerSecond:     rate.Limit(20),
				BurstSize:             20,
				MaxConcurrentRequests: 20,
				MaxRetries:            3,
				MaxRedirects:          5,
			},
			shouldUseDefault: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// If there's no config, ensure viper falls back to defaults
			viper.SetConfigFile("") // This disables reading config from file

			// Load the configuration
			LoadConfig()

			// Check if the config matches the expected values
			got := AppConfig
			want := tt.expectedConfig
			if got != want {
				t.Errorf("Expected config: %+v, but got: %+v", want, got)
			}
		})
	}
}
