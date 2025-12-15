// Package main is the entry point for the application.
package main

import (
	"golang-connect-marketplace/config"
	"golang-connect-marketplace/internal/auth/http/handlers"
	"golang-connect-marketplace/internal/auth/http/routes"
	"golang-connect-marketplace/internal/service"
	"golang-connect-marketplace/pkg/middleware"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

	e := echo.New()
	e.Use(middleware.RequestLogger(logger))

	setupAuth(e)

	err = e.Start(cfg.APIConfig.HTTPAddress)
	if err != nil {
		log.Error(err)
	}
}

func setupAuth(e *echo.Echo) {
	svc := service.NewService()
	hndl := handlers.NewHandler(svc)
	routes.RegisterRoutes(e, hndl)
}
