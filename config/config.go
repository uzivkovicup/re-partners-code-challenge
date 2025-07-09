package config

import (
	"errors"
	"fmt"
	"go-pack-calculator/constants"

	"github.com/spf13/viper"
)

// Config struct holds all configuration values
type Config struct {
	Environment        constants.AppEnv
	Protocol           string
	Host               string
	Port               int
	PostgresDBHost     string
	PostgresDBPort     int
	PostgresDBUser     string
	PostgresDBPassword string
	PostgresDBName     string
	PostgresDBSSLMode  string
}

// LoadConfig loads the configuration from environment variables and .env file
func Load() (*Config, error) {

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

	// Unmarshal config into struct
	config := &Config{
		Environment:        constants.AppEnv(viper.GetString("ENVIRONMENT")),
		Protocol:           viper.GetString("PROTOCOL"),
		Host:               viper.GetString("HOST"),
		Port:               viper.GetInt("PORT"),
		PostgresDBHost:     viper.GetString("POSTGRES_DB_HOST"),
		PostgresDBPort:     viper.GetInt("POSTGRES_DB_PORT"),
		PostgresDBUser:     viper.GetString("POSTGRES_DB_USER"),
		PostgresDBPassword: viper.GetString("POSTGRES_DB_PASSWORD"),
		PostgresDBName:     viper.GetString("POSTGRES_DB_NAME"),
		PostgresDBSSLMode:  viper.GetString("POSTGRES_DB_SSLMODE"),
	}

	return config, nil
}
