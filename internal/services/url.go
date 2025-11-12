package services

import (
	"url-shortener/internal/config"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

type ProcessingURL struct {
	repo storage.DbStore
	cfg  *config.Config
}

func NewProcessingURL(repo storage.DbStore, cfg *config.Config) *ProcessingURL {
	return &ProcessingURL{repo: repo, cfg: cfg}
}

func (s *ProcessingURL) SaveURL(urlToSave, alias string) (string, error) {
	if alias == "" {
		alias = random.NewRandomString(s.cfg.HTTPServer.AliasLength)
	}
	err := s.repo.SaveURL(urlToSave, alias)
	if err != nil {
		return "", err
	}
	return alias, nil
}

func (s *ProcessingURL) DeleteURL(alias string) error {
	return s.repo.DeleteURL(alias)
}

func (s *ProcessingURL) Redirect(alias string) (string, error) {
	return s.repo.GetURL(alias)
}
