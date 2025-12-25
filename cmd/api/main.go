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
	"golang-connect-marketplace/internal/marketplace/paymentproviders"
	marketRepos "golang-connect-marketplace/internal/marketplace/repos"
	marketSvc "golang-connect-marketplace/internal/marketplace/services"
	localStorage "golang-connect-marketplace/internal/marketplace/storage/local"
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
	e.Static(cfg.StorageConfig.UploadDir, cfg.StorageConfig.UploadDir)

	e.Use(middleware.RequestLogger(logger))
	e.Use(echoMiddleware.BodyLimit(cfg.APIConfig.MaxPayloadSize))

	authSvc := setupAuth(e, db, &cfg.AuthConfig)
	listingsRepo := setupListings(e, db, authSvc, &cfg.StorageConfig)
	setupPayments(e, db, authSvc, listingsRepo, &cfg.PaymentsConfig)

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

func setupListings( //nolint:ireturn
	e *echo.Echo,
	db *sqlx.DB,
	authSvc *authSvc.Service,
	cfg *config.StorageConfig,
) marketRepos.ListingsRepo {
	repo := marketRepos.NewListingsRepo(db)
	storage := localStorage.NewLocalStorage(cfg.UploadDir)
	svc := marketSvc.NewListingsService(repo, storage, cfg)
	hndl := marketHndl.NewListingsHandler(svc)
	marketRoutes.RegisterListingsRoutes(e, hndl, authSvc)

	return repo
}

func setupPayments(
	e *echo.Echo,
	db *sqlx.DB,
	authSvc *authSvc.Service,
	listingsRepo marketRepos.ListingsRepo,
	cfg *config.PaymentsConfig,
) {
	repo := marketRepos.NewPaymentsRepo(db)
	paymentProvider := paymentproviders.NewStripePaymentProvider(
		cfg.StripeSecretKey,
		cfg.StripeWebhookSecret,
	)
	svc := marketSvc.NewPaymentsService(paymentProvider, repo, listingsRepo)
	hndl := marketHndl.NewPaymentsHandler(svc)
	marketRoutes.RegisterPaymentsRoutes(e, hndl, authSvc)
}
