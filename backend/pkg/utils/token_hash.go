package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// HashToken returns a hex-encoded SHA-256 digest of the token.
func HashToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

// LooksLikeTokenHash reports whether s is a 64-char hex digest (stored API token hash).
func LooksLikeTokenHash(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, c := range strings.ToLower(s) {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}
