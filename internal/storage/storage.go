package storage

import (
	"fmt"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
	"url-shortener/internal/storage/sqlite"
)

type DbStore interface {
	SaveURL(urlToSave string, alias string) error
	GetURL(alias string) (string, error)
	DeleteURL(alias string) error
	Close()
}

func NewStorage(cfg *config.Config) (DbStore, error) {
	switch cfg.Database.Type {
	case "sqlite":
		return sqlite.New(cfg.DSN)
	case "postgres":
		return postgres.New(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Database.Type)
	}
}
