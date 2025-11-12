package handlers

import (
	"log/slog"
	v1 "url-shortener/internal/http-server/handlers/v1"
	mwLogger "url-shortener/internal/http-server/middleware/logger"

	_ "url-shortener/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(log *slog.Logger, handlers *v1.HTTPHandlers) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", handlers.SaveHandler())
	router.Delete("/url/{alias}", handlers.DeleteHandler())
	router.Get("/{alias}", handlers.RedirectHandler())
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return router
}
