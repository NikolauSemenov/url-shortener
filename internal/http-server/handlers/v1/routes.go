package v1

import (
	"net/http"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/ports"

	_ "url-shortener/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handlers interface {
	SaveHandler() http.HandlerFunc
	DeleteHandler() http.HandlerFunc
	RedirectHandler() http.HandlerFunc
}

func RegisterRoutes(log ports.Logger, router *chi.Mux, handlers Handlers) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/api/v1/url/save", handlers.SaveHandler())
	router.Delete("/api/v1/url/{alias}", handlers.DeleteHandler())
	router.Get("/api/v1/{alias}", handlers.RedirectHandler())
	router.Get("/api/v1/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/swagger/doc.json"),
	),
	)
}
