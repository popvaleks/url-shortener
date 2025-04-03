package redirect

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

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

// New
// @Summary Redirect by alias
// @Description Redirects to the original URL associated with the provided alias
// @Tags url
// @Param alias path string true "Alias for the URL to redirect"
// @Success 302 "Redirects to the original URL"
// @Failure 400 {object} resp.Response "Alias is missing"
// @Failure 404 {object} resp.Response "URL not found for the provided alias"
// @Failure 500 {object} resp.Response "Internal server error"
// @Router /{alias} [get]
func New(log *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

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

		rUrl, err := urlGetter.GetUrl(alias)
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

		log.Info("success get url", slog.String("res_url", rUrl))

		http.Redirect(w, r, rUrl, http.StatusFound)
	}
}
