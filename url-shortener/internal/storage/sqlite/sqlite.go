package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // init sqlite3 driver
)

type StorageSqlite struct {
	db *sql.DB
}

func New(storagePath string) (*StorageSqlite, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &StorageSqlite{db: db}, nil
}

func (s *StorageSqlite) Close() {
	s.db.Close()
}

func (s *StorageSqlite) SaveURL(urlToSave string, alias string) error {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO urls (url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StorageSqlite) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT urls FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resURL, nil

}

func (s *StorageSqlite) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"
	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}
