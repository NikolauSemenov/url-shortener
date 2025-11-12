package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers"
	v1 "url-shortener/internal/http-server/handlers/v1"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	web       *http.Server
	configApp *config.Config
	store     storage.DbStore
	log       *slog.Logger
}

func NewApp() (*App, error) {
	configApp := config.GetConfig()
	if configApp == nil {
		return nil, errors.New("configApp is nil")
	}

	log := setupLogger(configApp.Env)

	store, err := storage.NewStorage(configApp)
	if err != nil {
		return nil, err
	}

	srv := services.New(store, configApp)

	handlerGroupV1 := v1.NewHTTPHandlers(log, srv.URL)
	router := handlers.RegisterRoutes(log, handlerGroupV1)
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

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
