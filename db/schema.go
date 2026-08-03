package db

import (
	"database/sql"
	"strings"
)

func schema(db *sql.DB) error {
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;
		PRAGMA journal_mode = WAL;

		CREATE TABLE IF NOT EXISTS scripts (
			name    TEXT PRIMARY KEY,
			cmd     TEXT NOT NULL,
			args    TEXT NOT NULL DEFAULT '[]',
			runner  TEXT,
			workspace TEXT NOT NULL DEFAULT ''
		);

		CREATE TABLE IF NOT EXISTS chains (
			id      INTEGER PRIMARY KEY AUTOINCREMENT,
			name    TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS chain_steps (
			chain_id    INTEGER NOT NULL REFERENCES chains(id) ON DELETE CASCADE,
			script_name TEXT    NOT NULL REFERENCES scripts(name) ON DELETE RESTRICT,
			seq         INTEGER NOT NULL,
			pipe_after  INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (chain_id, seq)
		);	
	`)
	if err != nil {
		return err
	}

	// Migrate databases created before workspace support.
	_, err = db.Exec("ALTER TABLE scripts ADD COLUMN workspace TEXT NOT NULL DEFAULT ''")
	if err != nil && !strings.Contains(err.Error(), "duplicate column name: workspace") {
		return err
	}
	_, err = db.Exec("ALTER TABLE chain_steps ADD COLUMN pipe_after INTEGER NOT NULL DEFAULT 0")
	if err != nil && !strings.Contains(err.Error(), "duplicate column name: pipe_after") {
		return err
	}
	return nil
}
