package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/popvaleks/url-shortener/internal/storage"
	"log/slog"
	"net/http"

	resp "github.com/popvaleks/url-shortener/internal/lib/api/response"
	rand "github.com/popvaleks/url-shortener/internal/lib/utils/random"
)

type UrlSaver interface {
	SaveUrl(inputUrl string, alias string) (int64, error)
}

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 8

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

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

		alias := req.Alias
		if alias == "" {
			alias = rand.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveUrl(req.Url, alias)
		if errors.Is(err, storage.ErrUrlExists) {
			log.Info("url already exists", slog.String("url", req.Url))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to save url", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to save url"))

			return
		}

		log.Info("success save url", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Alias:    alias,
		})
	}
}
