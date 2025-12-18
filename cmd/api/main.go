// Package main is the entry point for the application.
package main

import (
	"golang-connect-marketplace/config"
	"golang-connect-marketplace/internal/auth/http/handlers"
	"golang-connect-marketplace/internal/auth/http/routes"
	authRepo "golang-connect-marketplace/internal/auth/repo"
	authSvc "golang-connect-marketplace/internal/auth/service"
	"golang-connect-marketplace/pkg/middleware"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	_ "github.com/lib/pq"
)

func main() {
	var cfg config.AppConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Panic("failed to load config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	db, err := sqlx.Connect(cfg.DBConfig.DriverName, cfg.DBConfig.URI)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Error("failed to ping database", "error", err)
	}

	logger.Info("connected to database")

	e := echo.New()
	e.Use(middleware.RequestLogger(logger))

	setupAuth(e, db, &cfg.AuthConfig)

	err = e.Start(cfg.APIConfig.HTTPAddress)
	if err != nil {
		log.Error(err)
	}
}

func setupAuth(e *echo.Echo, db *sqlx.DB, cfg *config.AuthConfig) {
	repo := authRepo.New(db)
	svc := authSvc.New(repo, cfg)
	hndl := handlers.NewHandler(svc)
	routes.RegisterRoutes(e, hndl)
}
