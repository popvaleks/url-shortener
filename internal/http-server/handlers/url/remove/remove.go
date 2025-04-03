package remove

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"

	resp "github.com/popvaleks/url-shortener/internal/lib/api/response"
	"github.com/popvaleks/url-shortener/internal/storage"
)

type UrlRemover interface {
	DeleteUrl(alias string) error
}

// Request represents URL deletion request
// @Description Request to delete a short URL
type Request struct {
	Alias string `json:"alias" validate:"required,alias"`
}

// Response represents URL deletion response
// @Description Success response for URL deletion
// swagger:model
type Response struct {
	resp.Response
}

// New
// @Summary Delete URL by alias
// @Description Deletes a short URL by its alias
// @Tags url
// @Param alias path string true "Alias of the URL to delete"
// @Success 200 {object} Response
// @Failure 400 {object} resp.Response "Alias is missing"
// @Failure 404 {object} resp.Response "URL not found"
// @Failure 500 {object} resp.Response "Internal server error"
// @Router /{alias} [delete]
func New(log *slog.Logger, urlRemover UrlRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.remove.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias not allowed")

			render.JSON(w, r, resp.Error("alias not allowed"))

			return
		}

		err := urlRemover.DeleteUrl(alias)
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("url not found")

			render.JSON(w, r, resp.Error("url not found"))

			return
		}

		if err != nil {
			log.Error("internal server error")

			render.JSON(w, r, resp.Error("internal server error"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
