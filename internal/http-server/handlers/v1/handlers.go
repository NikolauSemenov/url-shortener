package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/http-server/handlers/model"
	"url-shortener/internal/lib/api/errorsApp"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/ports"
	"url-shortener/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type HTTPHandlers struct {
	log ports.Logger
	svc services.URLService
}

func NewHTTPHandlers(log ports.Logger, svc services.URLService) *HTTPHandlers {
	return &HTTPHandlers{log: log, svc: svc}
}

// SaveHandler сохраняет URL с optional alias
// @Summary Сохраняет URL
// @Description Сохраняет оригинальный URL и возвращает alias. Если alias не передан, генерируется случайный.
// @Tags URL
// @Param request body model.Request true "Данные запроса"
// @Success 201 {string} string  "URL успешно сохранён"
// @Failure 400 {string} string "Некорректный JSON"
// @Failure 422 {string} string "Ошибка валидации запроса"
// @Failure 500 {string} string "Ошибка сохранения URL"
// @Router /api/v1/url/save [post]
func (h *HTTPHandlers) SaveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := h.log.With(
			slog.String("op", "SaveHandler"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req model.Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("failed to decode", "error", err)
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, "invalid json")
			return
		}

		if err := validator.New().Struct(req); err != nil {
			logger.Error("invalid request", "error", err)
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, err.Error())
			return
		}

		alias, err := h.svc.SaveURL(req.URL, req.Alias)
		if err != nil {
			if errors.Is(err, errorsApp.ErrUrlAlreadyExists) {
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, "url already exists")
				return
			}
			logger.Error("failed to save url", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "cannot save url")
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, model.Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}

// DeleteHandler удаляет URL по alias
// @Summary Удаляет URL
// @Description Удаляет URL по alias. Возвращает статус 204 No Content в случае успеха.
// @Tags URL
// @Param alias path string true "Alias для удаления"
// @Success 204 "URL успешно удалён"
// @Failure 422 {string} string "Alias не указан"
// @Failure 500 {string} string "Ошибка удаления URL"
// @Router /api/v1/url/{alias} [delete]
func (h *HTTPHandlers) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := h.log.With(
			slog.String("op", "DeleteHandler"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			render.JSON(w, r, "alias is required")
			return
		}

		if err := h.svc.DeleteURL(alias); err != nil {
			logger.Error("failed to delete", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "server error")
			return
		}

		render.Status(r, http.StatusNoContent)
	}
}

// RedirectHandler обрабатывает редирект по alias
// @Summary Перенаправляет на оригинальный URL по alias
// @Description Получает alias из URL, ищет соответствующий оригинальный URL и выполняет редирект (HTTP 302).
// @Tags URL
// @Param alias path string true "Alias для редиректа"
// @Success 302 {string} string "Redirected to original URL"
// @Router /api/v1/{alias} [get]
func (h *HTTPHandlers) RedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := h.log.With(
			slog.String("op", "RedirectHandler"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, "param alias is required")
			return
		}

		url, err := h.svc.Redirect(alias)
		if err != nil {
			logger.Error("error get url from alias", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, "server error")
			return
		}

		if url == "" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, "url not found")
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
