package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var weakJWTSecrets = []string{
	"your-secret-key",
	"your-secret-key-change-in-production",
}

type Config struct {
	Environment string
	Port        string
	GinMode     string

	Database    DatabaseConfig
	JWT         JWTConfig
	Storage     StorageConfig
	Email       EmailConfig
	CORS        CORSConfig
	RateLimit   RateLimitConfig
	AppURL      string
	Stripe      StripeConfig
	OAuth       OAuthConfig
	InternalKey     string
	TrustedProxies  []string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	RedirectBaseURL    string
}

type DatabaseConfig struct {
	URL      string // Full database URL (takes precedence if provided)
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

type StorageConfig struct {
	Type  string
	MinIO MinIOConfig
	S3    S3Config
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type S3Config struct {
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
}

type EmailConfig struct {
	ResendAPIKey string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type RateLimitConfig struct {
	Enabled      bool
	Requests     int // authenticated / default
	Window       time.Duration
	ListLimit    int // public list endpoints
	DetailLimit  int // public detail endpoints
	TrustedLimit int // SSR with internal key
}

// StripeConfig holds Stripe billing settings.
//
// Plan mapping:
//   - Job Boost (one-time, €49)   -> PriceJobBoost
//   - Startup Pro (subscription, €99/mo) -> PriceStartupPro
type StripeConfig struct {
	SecretKey       string
	WebhookSecret   string
	PriceJobBoost   string
	PriceStartupPro string
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if it doesn't)
	_ = godotenv.Load()

	config := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),

		Database: DatabaseConfig{
			URL:      getEnv("DATABASE_URL", ""),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "startup_board"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},

		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "your-secret-key"),
			Expiration:        parseDuration(getEnv("JWT_EXPIRATION", "24h")),
			RefreshExpiration: parseDuration(getEnv("JWT_REFRESH_EXPIRATION", "168h")),
		},

		Storage: StorageConfig{
			Type: getEnv("STORAGE_TYPE", "minio"),
			MinIO: MinIOConfig{
				Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
				AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
				SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
				Bucket:    getEnv("MINIO_BUCKET", "startup-board-uploads"),
				UseSSL:    getEnvBool("MINIO_USE_SSL", false),
			},
			S3: S3Config{
				Bucket:    getEnv("S3_BUCKET", ""),
				Region:    getEnv("S3_REGION", "us-east-1"),
				AccessKey: getEnv("AWS_ACCESS_KEY_ID", ""),
				SecretKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			},
		},

		Email: EmailConfig{
			ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		},

		CORS: CORSConfig{
			AllowedOrigins: parseStringSlice(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		},

		RateLimit: RateLimitConfig{
			Enabled:      getEnvBool("RATE_LIMIT_ENABLED", true),
			Requests:     getEnvInt("RATE_LIMIT_REQUESTS", 100),
			Window:       parseDuration(getEnv("RATE_LIMIT_WINDOW", "1m")),
			ListLimit:    getEnvInt("RATE_LIMIT_LIST_REQUESTS", 30),
			DetailLimit:  getEnvInt("RATE_LIMIT_DETAIL_REQUESTS", 60),
			TrustedLimit: getEnvInt("RATE_LIMIT_TRUSTED_REQUESTS", 300),
		},

		AppURL: getEnv("APP_URL", "http://localhost:3000"),

		InternalKey: getEnv("API_INTERNAL_KEY", ""),

		Stripe: StripeConfig{
			SecretKey:       getEnv("STRIPE_SECRET_KEY", ""),
			WebhookSecret:   getEnv("STRIPE_WEBHOOK_SECRET", ""),
			PriceJobBoost:   getEnv("STRIPE_PRICE_JOB_BOOST", ""),
			PriceStartupPro: getEnv("STRIPE_PRICE_STARTUP_PRO", ""),
		},

		OAuth: OAuthConfig{
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectBaseURL:    getEnv("OAUTH_REDIRECT_BASE_URL", "http://localhost:8080/api/v1/auth/oauth"),
		},

		TrustedProxies: parseStringSlice(getEnv("TRUSTED_PROXIES", "")),
	}

	if err := validateJWTSecret(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateJWTSecret(cfg *Config) error {
	isProd := cfg.Environment == "production" || cfg.Environment == "prod" || cfg.GinMode == "release"
	if !isProd {
		return nil
	}

	secret := cfg.JWT.Secret
	if secret == "" || len(secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters in production")
	}
	for _, weak := range weakJWTSecrets {
		if secret == weak {
			return fmt.Errorf("JWT_SECRET must not use a default value in production")
		}
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour // Default
	}
	return d
}

func parseStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	for _, item := range splitString(s, ",") {
		if trimmed := trimString(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	result := []string{}
	current := ""
	for _, char := range s {
		if string(char) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trimString(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
