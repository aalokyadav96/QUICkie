package initdb

import (
	"database/sql"
	"fmt"
)

// initDB opens (or creates) a SQLite database and ensures
// that the required table is created.
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Create the events table if it does not exist.
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entity_type TEXT,
		action TEXT,
		entity_id TEXT,
		item_id TEXT,
		item_type TEXT,
		additional_info TEXT
	);`
	if _, err = db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}
