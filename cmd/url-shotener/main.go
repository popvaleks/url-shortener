package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"

	"github.com/popvaleks/url-shortener/internal/config"
	"github.com/popvaleks/url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// CONFIG_PATH=config/local.yaml
func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("ENVIRONMENT", cfg.Env))
	log.Debug("debug start message")

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("error creating storage", err)
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

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
