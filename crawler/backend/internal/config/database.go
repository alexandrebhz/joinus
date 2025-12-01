package config

import (
	"fmt"
	"strings"

	"github.com/startup-job-board/crawler/internal/infrastructure/persistence/gorm_model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase creates a new database connection
func NewDatabase(databaseURL string) (*gorm.DB, error) {
	// Try to connect to the target database
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		// If connection fails, try to create the database
		if strings.Contains(err.Error(), "does not exist") {
			// Extract database name from URL
			dbName := extractDatabaseName(databaseURL)
			if dbName != "" {
				// Connect to postgres database to create the target database
				postgresURL := strings.Replace(databaseURL, "/"+dbName+"?", "/postgres?", 1)
				postgresURL = strings.Replace(postgresURL, "/"+dbName+"", "/postgres", 1)
				
				postgresDB, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
				if err != nil {
					return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
				}
				
				// Create the database if it doesn't exist
				sqlDB, err := postgresDB.DB()
				if err != nil {
					return nil, err
				}
				
				// Quote database name to handle special characters safely
				createDBQuery := fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)
				_, err = sqlDB.Exec(createDBQuery)
				if err != nil && !strings.Contains(err.Error(), "already exists") {
					sqlDB.Close()
					return nil, fmt.Errorf("failed to create database: %w", err)
				}
				
				sqlDB.Close()
				
				// Now try to connect to the newly created database
				db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
				if err != nil {
					return nil, fmt.Errorf("failed to connect after creating database: %w", err)
				}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&gorm_model.CrawlSiteModel{},
		&gorm_model.CrawledJobModel{},
		&gorm_model.CrawlLogModel{},
	); err != nil {
		return nil, err
	}

	return db, nil
}

// extractDatabaseName extracts the database name from a PostgreSQL connection URL
func extractDatabaseName(databaseURL string) string {
	// Parse URL format: postgres://user:pass@host:port/dbname?params
	parts := strings.Split(databaseURL, "/")
	if len(parts) < 4 {
		return ""
	}
	
	dbPart := parts[len(parts)-1]
	// Remove query parameters
	dbPart = strings.Split(dbPart, "?")[0]
	
	return dbPart
}

