package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"go-pack-calculator/db/migrations"
)

// Migration represents a database migration
type Migration struct {
	Version   string    `gorm:"primaryKey"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

// RunMigrations runs migrations
func RunMigrations(db *gorm.DB) error {
	// Create migrations table if not exists
	err := db.AutoMigrate(&Migration{})
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	var appliedMigrations []Migration
	if err := db.Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	applied := make(map[string]struct{})
	for _, m := range appliedMigrations {
		applied[m.Version] = struct{}{}
	}

	// Get all registered migrations
	migrations := migrations.List()
	for _, m := range migrations {
		if _, ok := applied[m.Version]; ok {
			log.Printf("Migration %s already applied\n", m.Version)

			continue
		}

		// Start a transaction for the migration
		tx := db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to start transaction for migration %s: %w", m.Version, tx.Error)
		}

		if err := m.Up(tx); err != nil {
			tx.Rollback()

			return fmt.Errorf("migration %s failed: %w", m.Version, err)
		}

		// Record the migration
		if err := tx.Create(&Migration{Version: m.Version}).Error; err != nil {
			tx.Rollback()

			return fmt.Errorf("failed to record migration %s: %w", m.Version, err)
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", m.Version, err)
		}

		log.Printf("Applied migration: %s\n", m.Version)
	}

	return nil
}
