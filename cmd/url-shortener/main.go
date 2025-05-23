package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/popvaleks/url-shortener/docs"
	"github.com/popvaleks/url-shortener/internal/config"
	"github.com/popvaleks/url-shortener/internal/http-server/handlers/url/getAllUrls"
	"github.com/popvaleks/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/popvaleks/url-shortener/internal/http-server/handlers/url/remove"
	"github.com/popvaleks/url-shortener/internal/http-server/handlers/url/save"
	"github.com/popvaleks/url-shortener/internal/http-server/handlers/url/updateUrl"
	mwLogger "github.com/popvaleks/url-shortener/internal/http-server/middleware/logger"
	"github.com/popvaleks/url-shortener/internal/storage/sqlite"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Url shortener
// @version 1.0
// @description Shortener service
// @host localhost:8080
// @BasePath /
func main() {
	// CONFIG_PATH=config/local.yaml
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("ENVIRONMENT", cfg.Env))
	log.Debug("debug start message")

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("error creating storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer) // anti panic
	router.Use(middleware.URLFormat) // routing

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	router.Delete("/{alias}", remove.New(log, storage))
	router.Get("/url", getAllUrls.New(log, storage))
	router.Patch("/{alias}", updateUrl.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("error starting server", slog.String("error", err.Error()))
	}

	log.Error("stopping server", slog.String("address", cfg.Address))
}

func setupLogger(env string) (log *slog.Logger) {
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
