package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // init sqlite3 driver
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w, error whith opening", op, err)
	}

	// https://www.dfkg.com/tyvubiunoipmo
	// https://my.sh/sqlite
	// protocol://domen/alias

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS indx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w, error with creating", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w, error with starting", op, err)
	}

	return &Storage{db: db}, nil

}
