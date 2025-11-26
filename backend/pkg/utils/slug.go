package utils

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
)

var (
	slugRegex = regexp.MustCompile("[^a-z0-9]+")
)

func GenerateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces and special chars with hyphens
	slug = slugRegex.ReplaceAllString(slug, "-")
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	return slug
}

func GenerateUniqueSlug(baseSlug string, checkExists func(string) bool) string {
	slug := baseSlug
	counter := 0

	for checkExists(slug) {
		// Generate random number suffix
		randomBytes := make([]byte, 4)
		rand.Read(randomBytes)
		randomNum := fmt.Sprintf("%x", randomBytes)[:8]
		slug = fmt.Sprintf("%s-%s", baseSlug, randomNum)
		counter++
		// Safety check to avoid infinite loop
		if counter > 100 {
			break
		}
	}

	return slug
}



