package valueobject

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrInvalidTokenPrefix = errors.New("invalid token prefix")
)

const (
	TokenPrefixStartup = "sb_" // startup board token
	TokenLength        = 32    // base64 encoded length
)

type APIToken struct {
	prefix string
	value  string
}

func NewAPIToken(prefix string) (APIToken, error) {
	if prefix == "" {
		prefix = TokenPrefixStartup
	}

	// Generate random bytes
	bytes := make([]byte, 24) // 24 bytes = 32 base64 chars
	if _, err := rand.Read(bytes); err != nil {
		return APIToken{}, err
	}

	// Encode to base64
	token := base64.URLEncoding.EncodeToString(bytes)
	// Remove padding
	token = strings.TrimRight(token, "=")

	return APIToken{
		prefix: prefix,
		value:  token,
	}, nil
}

func ParseAPIToken(token string) (APIToken, error) {
	parts := strings.Split(token, "_")
	if len(parts) < 2 {
		return APIToken{}, ErrInvalidTokenFormat
	}

	prefix := strings.Join(parts[:len(parts)-1], "_") + "_"
	value := parts[len(parts)-1]

	if prefix != TokenPrefixStartup {
		return APIToken{}, ErrInvalidTokenPrefix
	}

	return APIToken{
		prefix: prefix,
		value:  value,
	}, nil
}

func (t APIToken) String() string {
	return fmt.Sprintf("%s%s", t.prefix, t.value)
}

func (t APIToken) Value() string {
	return t.value
}

func (t APIToken) Prefix() string {
	return t.prefix
}



