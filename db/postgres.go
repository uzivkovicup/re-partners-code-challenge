package db

import (
	"fmt"
	"go-pack-calculator/constants"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresDB is the global PostgresDB instance
var PostgresDB *gorm.DB

// postgresDB struct
type PostgresDBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// NewPostgresDB now only initializes the config struct
func NewPostgresDB(host string, port int, user, password, database, sslMode string) *PostgresDBConfig {
	return &PostgresDBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  sslMode,
	}
}

// Connect to PostgresDB
func (p *PostgresDBConfig) Connect() error {
	if PostgresDB != nil {
		return nil
	}

	// Ensure DB exists before connecting
	err := EnsureDatabaseExists(p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode)
	if err != nil {
		return fmt.Errorf("failed to ensure database exists: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode,
	)

	var db *gorm.DB
	maxRetries := 5
	delay := 2 // seconds

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             constants.GormLoggerSlowThreshold,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Try to connect with retries
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})

		if err == nil {
			// Test connection
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					break // success
				}
			}
		}

		if i < maxRetries-1 {
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}

	if err != nil {
		return fmt.Errorf("could not connect to postgres after %d attempts: %w", maxRetries, err)
	}

	log.Printf("Successfully connected to database %s\n", p.Database)

	// Run migrations
	if err := RunMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	PostgresDB = db

	return nil
}

// Close PostgresDB
func (p *PostgresDBConfig) Close() error {
	if PostgresDB == nil {
		return nil
	}

	sqlDB, err := PostgresDB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
