package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func VerifyHash(input, hash string) bool {
	computedHash := HashString(input)
	return computedHash == hash
}



