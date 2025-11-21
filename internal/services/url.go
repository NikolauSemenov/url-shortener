package services

import (
	"errors"
	"log/slog"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/api/errorsApp"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/cache"
)

type ProcessingURL struct {
	repo        storage.DbStore
	cfg         *config.Config
	cacheClient cache.Cache
	log         *slog.Logger
}

func NewProcessingURL(repo storage.DbStore, cfg *config.Config, cacheClient cache.Cache, log *slog.Logger) *ProcessingURL {
	return &ProcessingURL{repo: repo, cfg: cfg, cacheClient: cacheClient, log: log}
}

func (s *ProcessingURL) SaveURL(urlToSave, alias string) (string, error) {
	if alias == "" {
		alias = random.NewRandomString(s.cfg.HTTPServer.AliasLength)
	} else {
		if err := s.repo.CheckExistsUrl(alias, urlToSave); err != nil && errors.Is(err, errorsApp.ErrUrlAlreadyExists) {
			return "", errorsApp.ErrUrlAlreadyExists
		}
	}
	err := s.repo.SaveURL(urlToSave, alias)
	if err != nil {
		return "", err
	}
	return alias, nil
}

func (s *ProcessingURL) DeleteURL(alias string) error {
	err := s.cacheClient.Invalidate(alias)
	if err != nil {
		s.log.Error("failed to invalidate cache url", "alias", alias, "err", err)
	}
	return s.repo.DeleteURL(alias)
}

func (s *ProcessingURL) Redirect(alias string) (string, error) {
	data, err := s.cacheClient.Get(alias)
	if err != nil {
		if errors.Is(err, errorsApp.ErrCacheMiss) {
			val, errRep := s.repo.GetURL(alias)
			if errRep != nil {
				return "", errRep
			}
			err = s.cacheClient.Set(alias, val, s.cfg.CacheTTL)
			if err != nil {
				s.log.Error("failed to create cache url", "alias", alias, "err", err)
			}
			return val, nil
		}
	}
	return data, nil
}
