package services

import (
	"log/slog"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/cache"
)

type URLService interface {
	SaveURL(urlToSave, alias string) (string, error)
	DeleteURL(alias string) error
	Redirect(alias string) (string, error)
}

type Services struct {
	URL URLService
	log *slog.Logger
}

func New(store storage.DbStore, cfg *config.Config, cacheClient cache.Cache, log *slog.Logger) *Services {
	return &Services{
		URL: NewProcessingURL(store, cfg, cacheClient, log),
	}
}
