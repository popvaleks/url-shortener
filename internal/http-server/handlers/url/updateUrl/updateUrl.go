package updateUrl

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/popvaleks/url-shortener/internal/storage"
	"log/slog"
	"net/http"

	resp "github.com/popvaleks/url-shortener/internal/lib/api/response"
)

type UrlEditer interface {
	UpdateUrl(url, alias string) (string, error)
}

// Request represents URL update request
// @Description Request to update original URL for existing alias
type Request struct {
	Url string `json:"url" validate:"required,url"`
}

// ResponseAlias represents alias in response
// @Description Contains updated alias information
type ResponseAlias struct {
	Alias string `json:"alias"`
}

// Response represents URL update response
// @Description Success response with updated alias
// swagger:model
type Response struct {
	resp.Response
	Result ResponseAlias `json:"result"`
}

// New
// @Summary Update URL by alias
// @Description Updates original URL for existing alias
// @Tags url
// @Accept  json
// @Produce  json
// @Param alias path string true "Alias to update"
// @Param input body Request true "New URL data"
// @Success 200 {object} Response
// @Failure 400 {object} resp.Response "Invalid request or validation error"
// @Failure 404 {object} resp.Response "Alias not found"
// @Failure 500 {object} resp.Response "Internal server error"
// @Router /{alias} [patch]
func New(log *slog.Logger, urlEditer UrlEditer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.updateUrl.New"

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

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("success decode req", slog.Any("request", req)) // can remove, debug only

		if err := validator.New().Struct(req); err != nil {
			var validatorErr validator.ValidationErrors
			errors.As(err, &validatorErr)

			log.Error("failed to validate request", slog.String("error", err.Error()))

			render.JSON(w, r, resp.ValidationError(validatorErr))

			return
		}

		sAlias, err := urlEditer.UpdateUrl(req.Url, alias)
		if errors.Is(err, storage.ErrAliasNotFound) {
			log.Info("alias not found", slog.String("url", req.Url))

			render.JSON(w, r, resp.Error("alias not found"))

			return
		}
		if err != nil {
			log.Error("failed to update url", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to update url"))

			return
		}

		log.Info("success update url", slog.String("alias", sAlias))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Result:   ResponseAlias{sAlias},
		})
	}
}
