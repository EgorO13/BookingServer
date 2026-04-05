package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewMigrationDriver(db *sql.DB) (*postgres.Postgres, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	pgDriver, ok := driver.(*postgres.Postgres)
	if !ok {
		return nil, fmt.Errorf("failed to cast driver to *postgres.Postgres")
	}
	return pgDriver, nil
}

func RunMigrations(driver *postgres.Postgres) error {
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
