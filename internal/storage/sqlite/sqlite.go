package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url/internal/storage"

	"github.com/mattn/go-sqlite3"
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

	// Как это работает?
	// Что такое stmt?
	// Что возвращает db.Prepare?
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

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w On prepare", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		//Как это работает? SQLite - Constraints 46:43
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w On exec", op, err)
	}

	id, err := res.LastInsertId()
	// ATTENTION DEBUG
	// fmt.Println(id)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
	// Как составляется запрос?
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statment: %w", op, err)
	}

	var resURL string
	// ОАОАОАОАООАОА Что это блять и как работает? 50:48
	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statment: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: deleting statment: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return fmt.Errorf("%s: %w On exec", op, err)
	}
	return nil
}
