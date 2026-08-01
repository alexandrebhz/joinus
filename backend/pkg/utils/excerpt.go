package utils

import "unicode/utf8"

// Excerpt truncates s to at most maxLen runes, appending "…" when truncated.
func Excerpt(s string, maxLen int) string {
	if maxLen <= 0 || s == "" {
		return s
	}
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	if maxLen == 1 {
		return "…"
	}
	return string(runes[:maxLen-1]) + "…"
}
