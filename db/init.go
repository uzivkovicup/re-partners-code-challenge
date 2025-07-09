package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// EnsureDatabaseExists ensures that the database exists
// This still uses the standard database/sql package because we need to connect to postgres database
// before our target database exists
func EnsureDatabaseExists(host string, port int, user, password, dbName, sslMode string) error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		host, port, user, password, sslMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return err
		}
		log.Printf("Database '%s' created.\n", dbName)
	}
	return nil
}
