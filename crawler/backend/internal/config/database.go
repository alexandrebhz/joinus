package config

import (
	"github.com/startup-job-board/crawler/internal/infrastructure/persistence/gorm_model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase creates a new database connection
func NewDatabase(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&gorm_model.CrawlSiteModel{},
		&gorm_model.CrawledJobModel{},
	); err != nil {
		return nil, err
	}

	return db, nil
}

