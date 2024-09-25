package utils

import "unicode"

// IsLetter is a wrapper over the `unicode.IsLetter()` function that returns true if all the runes in the string are letters (Unicode category "L").
func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
