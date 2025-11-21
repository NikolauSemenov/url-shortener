package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/config"
	v1 "url-shortener/internal/http-server/handlers/v1"
	"url-shortener/internal/infrastructure/logger"
	"url-shortener/internal/ports"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/cache"

	"github.com/go-chi/chi/v5"
)

type App struct {
	web       *http.Server
	configApp *config.Config
	store     storage.DbStore
	log       ports.Logger
}

func NewApp() (*App, error) {
	configApp := config.GetConfig()
	if configApp == nil {
		return nil, errors.New("configApp is nil")
	}

	log := logger.NewLogger(configApp.Env)

	store, err := storage.NewStorage(configApp)
	if err != nil {
		return nil, err
	}

	cacheClient, err := cache.NewCache(&configApp.CacheConfig)
	if err != nil {
		return nil, err
	}

	srv := services.New(store, configApp, cacheClient, log)

	handlerGroupV1 := v1.NewHTTPHandlers(log, srv.URL)

	router := chi.NewRouter()
	v1.RegisterRoutes(log, router, handlerGroupV1)
	return &App{
		web: &http.Server{
			Addr:         configApp.Address,
			Handler:      router,
			ReadTimeout:  configApp.HTTPServer.Timeout,
			WriteTimeout: configApp.HTTPServer.Timeout,
			IdleTimeout:  configApp.HTTPServer.IdleTimeout,
		},
		configApp: configApp,
		store:     store,
		log:       log,
	}, nil
}

func (a *App) Run() error {

	go func() {
		if err := a.web.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error(err.Error())
			os.Exit(1)
		}
	}()

	a.log.Info(fmt.Sprintf("Сервер запущен: %s", a.configApp.HTTPServer.Address))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	return a.Stop()
}

func (a *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.configApp.Timeout*time.Second)
	defer cancel()

	if err := a.web.Shutdown(ctx); err != nil {
		return err
	}
	a.store.Close()
	return nil
}
