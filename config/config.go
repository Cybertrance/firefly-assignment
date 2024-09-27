/*
Package config provides a centralized way to manage application configuration settings.
It uses Viper to load configuration values from both default settings and a "config.yaml"
file.
*/
package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// Config structure to hold the configuration
type Config struct {
	TopResults            int        `mapstructure:"top_results"`
	SourceURLFileName     string     `mapstructure:"source_url_filename"`
	WordBankURL           string     `mapstructure:"word_bank_url"`
	ContainerSelector     string     `mapstructure:"container_selector"`
	RequestsPerSecond     rate.Limit `mapstructure:"requests_per_second"`
	BurstSize             int        `mapstructure:"burst_size"`
	MaxConcurrentRequests int        `mapstructure:"max_concurrent_requests"`
	MaxRetries            int        `mapstructure:"max_retries"`
	MaxRedirects          int        `mapstructure:"max_redirects"`
}

var AppConfig Config

// LoadConfig loads configuration settings for the application.
// It sets default values for various parameters using Viper and attempts
// to read from a configuration file named "config.yaml" located in the current directory.
// If no config file is found, the function proceeds with the default values.
func LoadConfig() {
	// Set default values
	viper.SetDefault("top_results", 10)
	viper.SetDefault("source_url_filename", "endg-urls")
	viper.SetDefault("word_bank_url", "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt")
	viper.SetDefault("container_selector", ".caas-body")
	viper.SetDefault("requests_per_second", 20)
	viper.SetDefault("burst_size", 20)
	viper.SetDefault("max_concurrent_requests", 20)
	viper.SetDefault("max_retries", 3)
	viper.SetDefault("max_redirects", 5)

	// Configuration file settings
	viper.SetConfigName("config") // Config file name (without extension)
	viper.SetConfigType("yaml")   // Config file type
	viper.AddConfigPath(".")      // Look for config in the current directory

	// Read the config file if available; otherwise, continue with defaults
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("[INFO] - no configuration file found, using default values)")
	}

	// Unmarshal the config into AppConfig struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}
