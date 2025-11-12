package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoragePG struct {
	db *pgxpool.Pool
}

func New(storagePath string) (*StoragePG, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создание пула подключений
	pool, err := pgxpool.New(ctx, storagePath)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &StoragePG{db: pool}, nil
}

func (s *StoragePG) Close() {
	s.db.Close()
}

func (s *StoragePG) GetURL(alias string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var url string
	query := `SELECT original_url FROM urls WHERE alias = $1`
	err := s.db.QueryRow(ctx, query, alias).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

		}
		return "", nil
	}
	return url, nil
}

func (s *StoragePG) SaveURL(urlToSave string, alias string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO urls (alias, original_url) VALUES ($1, $2)`
	_, err := s.db.Exec(ctx, query, alias, urlToSave)
	return err
}

func (s *StoragePG) DeleteURL(alias string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM urls WHERE alias = $1`
	_, err := s.db.Exec(ctx, query, alias)
	if err != nil {
		return err
	}
	return nil
}
