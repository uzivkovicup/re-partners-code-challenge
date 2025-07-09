package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config struct holds all configuration values
type Config struct {
	Environment          string `mapstructure:"ENVIRONMENT"`
	Protocol             string `mapstructure:"PROTOCOL"`
	Host                 string `mapstructure:"HOST"`
	Port                 int    `mapstructure:"PORT"`
	PostgresDBHost       string `mapstructure:"POSTGRES_DB_HOST"`
	PostgresDBPort       int    `mapstructure:"POSTGRES_DB_PORT"`
	PostgresDBUser       string `mapstructure:"POSTGRES_DB_USER"`
	PostgresDBPassword   string `mapstructure:"POSTGRES_DB_PASSWORD"`
	PostgresDBName       string `mapstructure:"POSTGRES_DB_NAME"`
	PostgresDBSSLMode    string `mapstructure:"POSTGRES_DB_SSLMODE"`
	JWTSecret            string `mapstructure:"JWT_SECRET"`
	JWTExpiration        int    `mapstructure:"JWT_EXPIRATION"`
	JWTRefreshExpiration int    `mapstructure:"JWT_REFRESH_EXPIRATION"`
}

// Default config instance
var Default *Config

// LoadConfig loads the configuration from environment variables and .env file
func LoadConfig(configPath string) (*Config, error) {

	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.AllowEmptyEnv(true)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	// Check if .env file exists and load it
	envFile := filepath.Join(configPath, ".env")
	if _, err := os.Stat(envFile); err == nil {
		log.Printf("Loading configuration from %s", envFile)
		// Load .env file
		viper.SetConfigFile(envFile)
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		log.Printf(".env file not found at %s, using environment variables only", envFile)
	}

	// Unmarshal config into struct
	config := &Config{
		Environment:          viper.GetString("ENVIRONMENT"),
		Protocol:             viper.GetString("PROTOCOL"),
		Host:                 viper.GetString("HOST"),
		Port:                 viper.GetInt("PORT"),
		PostgresDBHost:       viper.GetString("POSTGRES_DB_HOST"),
		PostgresDBPort:       viper.GetInt("POSTGRES_DB_PORT"),
		PostgresDBUser:       viper.GetString("POSTGRES_DB_USER"),
		PostgresDBPassword:   viper.GetString("POSTGRES_DB_PASSWORD"),
		PostgresDBName:       viper.GetString("POSTGRES_DB_NAME"),
		PostgresDBSSLMode:    viper.GetString("POSTGRES_DB_SSLMODE"),
		JWTSecret:            viper.GetString("JWT_SECRET"),
		JWTExpiration:        viper.GetInt("JWT_EXPIRATION"),
		JWTRefreshExpiration: viper.GetInt("JWT_REFRESH_EXPIRATION"),
	}

	// Set as default config
	Default = config

	return config, nil
}

// Load is a backward-compatible wrapper for LoadConfig
func Load(filePath string) (*Config, error) {
	// Extract directory from file path
	dir := filepath.Dir(filePath)
	return LoadConfig(dir)
}
