package database

import (
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(pool *pgxpool.Pool, fs embed.FS) error {
	db := stdlib.OpenDBFromPool(pool) // Converts pool to *sql.DB
	goose.SetBaseFS(fs) // Tells goose where the files are

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("Failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("Failed to run migrations: %w", err)
	}

	return nil
}