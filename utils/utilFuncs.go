// Package utils provides utility functions for common operations used throughout the application.
package utils

import "unicode"

// IsLetter checks whether all characters in the provided string are letters.
// It returns true if all runes in the string belong to the Unicode letter category ("L").
//
// Parameters:
//   - s: The string to check.
//
// Returns:
//   - bool: True if all characters in the string are letters, false otherwise.
func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
