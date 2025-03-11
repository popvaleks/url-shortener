package getAllUrls

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"

	resp "github.com/popvaleks/url-shortener/internal/lib/api/response"
)

type AllUrlGetter interface {
	GetAllUrls() (map[string]string, error)
}

type Response struct {
	resp.Response
	Result map[string]string `json:"result"`
}

func New(log *slog.Logger, allUrlGetter AllUrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getAllUrls.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		rUrlMap, err := allUrlGetter.GetAllUrls()
		if err != nil {
			log.Error("internal server error")

			render.JSON(w, r, resp.Error("internal server error"))

			return
		}

		log.Info("success get all urls", slog.Group("res_urlMap", rUrlMap))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Result:   rUrlMap,
		})
	}
}
