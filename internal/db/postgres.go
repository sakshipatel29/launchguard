package db

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectPostgres() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = "postgres://launchguard:launchguard@localhost:5433/launchguard?sslmode=disable"
	}

	database, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}

func RunMigrations(database *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS feature_flags (
		id UUID PRIMARY KEY,
		name TEXT NOT NULL,
		flag_key TEXT NOT NULL,
		description TEXT,
		enabled BOOLEAN NOT NULL DEFAULT false,
		rollout_percentage INT NOT NULL CHECK (rollout_percentage >= 0 AND rollout_percentage <= 100),
		environment TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL,
		UNIQUE(flag_key, environment)
	);
	`

	_, err := database.Exec(query)
	return err
}
