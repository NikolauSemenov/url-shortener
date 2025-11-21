package services

import (
	"url-shortener/internal/config"
	"url-shortener/internal/ports"
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
	log *ports.Logger
}

func New(store storage.DbStore, cfg *config.Config, cacheClient cache.Cache, log ports.Logger) *Services {
	return &Services{
		URL: NewProcessingURL(store, cfg, cacheClient, log),
	}
}
