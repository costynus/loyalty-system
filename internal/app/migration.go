package app

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func migration(pgURL, migDir string) error {
	db, err := goose.OpenDBWithDriver("postgres", pgURL)
	if err != nil {
		return fmt.Errorf("migration - goose.OpenDBWithDriver: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, migDir); err != nil {
		fmt.Println(err)
		return fmt.Errorf("migration - goose.Up: %w", err)
	}
	return nil
}
