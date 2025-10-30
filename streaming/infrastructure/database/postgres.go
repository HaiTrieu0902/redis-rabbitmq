package database

import (
	"database/sql"
	"fmt"
	"go-api-streaming/infrastructure/config"

	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.GetDSN()
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	fmt.Println("✓ Database connection established")
	return db, nil
}
