package services

import (
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
)

type URLService interface {
	SaveURL(urlToSave, alias string) (string, error)
	DeleteURL(alias string) error
	Redirect(alias string) (string, error)
}

type Services struct {
	URL URLService
}

func New(store storage.DbStore, cfg *config.Config) *Services {
	return &Services{
		URL: NewProcessingURL(store, cfg),
	}
}
