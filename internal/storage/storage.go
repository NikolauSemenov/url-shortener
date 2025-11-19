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
	CheckExistsUrl(alias, urlToSave string) error
	Close()
}

func NewStorage(cfg *config.Config) (DbStore, error) {
	switch cfg.Database.DbType {
	case "sqlite":
		return sqlite.New(cfg.DbDsn)
	case "postgres":
		return postgres.New(cfg.DbDsn)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Database.DbType)
	}
}
