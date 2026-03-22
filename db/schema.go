package db

import "database/sql"

func schema(db *sql.DB) error {
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;
		PRAGMA journal_mode = WAL;

		CREATE TABLE IF NOT EXISTS scripts (
			name    TEXT PRIMARY KEY,
			cmd     TEXT NOT NULL,
			args    TEXT NOT NULL DEFAULT '[]',
			runner  TEXT
		);

		CREATE TABLE IF NOT EXISTS chains (
			id      INTEGER PRIMARY KEY AUTOINCREMENT,
			name    TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS chain_steps (
			chain_id    INTEGER NOT NULL REFERENCES chains(id) ON DELETE CASCADE,
			script_name TEXT    NOT NULL REFERENCES scripts(name) ON DELETE RESTRICT,
			seq         INTEGER NOT NULL,
			PRIMARY KEY (chain_id, seq)
		);	
	`)
	return err
}