package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/config"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
	platformmiddleware "github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/middleware"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/postgres"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/postgres/transactor"
	projecthttp "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/http"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/keygen"
	elementDB "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/repository/postgres/element"
	projectDB "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/repository/postgres/project"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/usecase/element"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/usecase/project"
	db "github.com/DaniilSintsov/interactive-onboarding/backend/migrations"
	"go.uber.org/zap"

	analyticshttp "github.com/DaniilSintsov/interactive-onboarding/backend/internal/analytics/http"
	analyticsrepository "github.com/DaniilSintsov/interactive-onboarding/backend/internal/analytics/repository"
	analyticsservice "github.com/DaniilSintsov/interactive-onboarding/backend/internal/analytics/service"
	pdfhttp "github.com/DaniilSintsov/interactive-onboarding/backend/internal/pdf_report/http"
)

func Run(logger *zap.Logger, cfg *config.Config) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer stop()

	poolCtx, cancelPC := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancelPC()

	dbPool, err := postgres.NewPool(poolCtx, cfg.ConstructPostgresURL())
	if err != nil {
		return err
	}
	defer dbPool.Close()

	migrationCtx, cancelMC := context.WithTimeout(ctx, time.Minute)
	defer cancelMC()
	if err = db.SetupPostgres(migrationCtx, dbPool); err != nil {
		return err
	}

	txManager := transactor.NewTransactor(dbPool)

	elementRepository := elementDB.NewRepository(dbPool)
	projectRepository := projectDB.NewRepository(dbPool)

	keyGenerator := keygen.NewGenerator()

	elementService := element.NewElementService(
		elementRepository,
		projectRepository,
		nil,
		txManager,
		logger,
	)

	projectService := project.NewProjectService(
		projectRepository,
		elementRepository,
		keyGenerator,
		elementService,
		txManager,
		logger,
	)

	projectHandler := projecthttp.NewHandler(
		elementService,
		projectService,
		logger,
	)

	analyticsRepo := analyticsrepository.NewAnalyticsRepository(dbPool)
	analyticsSvc := analyticsservice.NewAnalyticsService(analyticsRepo)
	analyticsHandler := analyticshttp.NewAnalyticsHandler(analyticsSvc)

	pdfHandler := pdfhttp.NewPDFHandler(analyticsSvc)

	server := httpserver.NewServer(
		cfg.HTTPConfig,
		platformmiddleware.CORS(cfg.HTTPConfig.AllowedOrigins),
	)

	server.RegisterRouteGroup(
		"/api/v1/",
		[]httpserver.Middleware{},
		projectHandler,
		analyticsHandler,
		pdfHandler,
	)

	if err = runHTTPServer(ctx, logger, cfg, server); err != nil {
		return fmt.Errorf("failed to run http server: %w", err)
	}

	return nil
}

func runHTTPServer(ctx context.Context, logger *zap.Logger, cfg *config.Config, server *httpserver.Server) error {
	serveErr := make(chan error, 1)

	go func() {
		logger.Info("http server started", zap.String("address", server.Address()))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- fmt.Errorf("http server listen error: %w", err)
			return
		}
		serveErr <- nil
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTPConfig.HTTPShutdownTime)
	defer cancel()

	logger.Info("shutting down http server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Warn("http server shutdown error", zap.Error(err))
		if closeErr := server.Close(); closeErr != nil && !errors.Is(closeErr, http.ErrServerClosed) {
			logger.Warn("http server forced close error", zap.Error(closeErr))
		}
		return fmt.Errorf("http server shutdown error: %w", err)
	}
	logger.Info("http server gracefully shutdown")

	return <-serveErr
}
