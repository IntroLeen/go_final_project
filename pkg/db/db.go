package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	if _, err := os.Stat(dbFile); err != nil {
	}

	d, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}
	if err := d.Ping(); err != nil {
		_ = d.Close()
		return fmt.Errorf("ping sqlite: %w", err)
	}

	if _, err := d.Exec(schema); err != nil {
		_ = d.Close()
		return fmt.Errorf("init schema: %w", err)
	}

	DB = d
	return nil
}

func Close() error {
	if DB != nil {
		err := DB.Close()
		DB = nil
		return err
	}
	return nil
}
