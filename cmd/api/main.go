// Package main is the entry point for the application.
package main

import (
	"golang-connect-marketplace/config"
	authHndl "golang-connect-marketplace/internal/auth/http/handlers"
	authRoutes "golang-connect-marketplace/internal/auth/http/routes"
	authRepo "golang-connect-marketplace/internal/auth/repo"
	authSvc "golang-connect-marketplace/internal/auth/service"
	marketHndl "golang-connect-marketplace/internal/marketplace/http/handlers"
	marketRoutes "golang-connect-marketplace/internal/marketplace/http/routes"
	marketRepos "golang-connect-marketplace/internal/marketplace/repos"
	marketSvc "golang-connect-marketplace/internal/marketplace/services"
	"golang-connect-marketplace/pkg/middleware"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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
	e.Use(echoMiddleware.BodyLimit(cfg.APIConfig.MaxPayloadSize))

	authSvc := setupAuth(e, db, &cfg.AuthConfig)
	setupListings(e, db, authSvc)

	err = e.Start(cfg.APIConfig.HTTPAddress)
	if err != nil {
		log.Error(err)
	}
}

func setupAuth(e *echo.Echo, db *sqlx.DB, cfg *config.AuthConfig) *authSvc.Service {
	repo := authRepo.New(db)
	svc := authSvc.New(repo, cfg)
	hndl := authHndl.NewHandler(svc)
	authRoutes.RegisterRoutes(e, hndl, svc)

	return svc
}

func setupListings(e *echo.Echo, db *sqlx.DB, authSvc *authSvc.Service) {
	repo := marketRepos.NewListingsRepo(db)
	svc := marketSvc.NewListingsService(repo)
	hndl := marketHndl.NewListingsHandler(svc)
	marketRoutes.RegisterRoutes(e, hndl, authSvc)
}
