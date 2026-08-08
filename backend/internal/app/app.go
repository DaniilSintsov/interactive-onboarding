package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	analyticshttp "github.com/DaniilSintsov/interactive-onboarding/backend/internal/analytics/http"
	analyticsrepository "github.com/DaniilSintsov/interactive-onboarding/backend/internal/analytics/repository"
	analyticsservice "github.com/DaniilSintsov/interactive-onboarding/backend/internal/analytics/service"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/config"
	pdfhttp "github.com/DaniilSintsov/interactive-onboarding/backend/internal/pdf_report/http"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
	platformmiddleware "github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/middleware"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/postgres"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/postgres/transactor"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/testtoken"
	projecthttp "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/http"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/keygen"
	elementDB "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/repository/postgres/element"
	projectDB "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/repository/postgres/project"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/usecase/element"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/usecase/project"
	runtimehttp "github.com/DaniilSintsov/interactive-onboarding/backend/internal/runtime/http/handlers"
	runtimescenario "github.com/DaniilSintsov/interactive-onboarding/backend/internal/runtime/repository/scenario"
	runtimesession "github.com/DaniilSintsov/interactive-onboarding/backend/internal/runtime/repository/session"
	runtimestep "github.com/DaniilSintsov/interactive-onboarding/backend/internal/runtime/repository/step"
	runtimetoken "github.com/DaniilSintsov/interactive-onboarding/backend/internal/runtime/repository/token"
	runtimeservice "github.com/DaniilSintsov/interactive-onboarding/backend/internal/runtime/service"
	scenariohttp "github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/http"
	scenarioDB "github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/repository/postgres/scenario"
	stepDB "github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/repository/postgres/step"
	testTokenDB "github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/repository/postgres/test_token"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/usecase/scenario"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/usecase/step"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/usecase/test_token"
	trackinghttp "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/http/handlers"
	trackingevent "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/repository/event"
	trackingscenario "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/repository/scenario"
	trackingsession "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/repository/session"
	trackingstep "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/repository/step"
	trackingservice "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/service"
	db "github.com/DaniilSintsov/interactive-onboarding/backend/migrations"
	"go.uber.org/zap"
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
	stepRepository := stepDB.NewRepository(dbPool, txManager)
	testTokenRepository := testTokenDB.NewRepository(dbPool)
	scenarioRepository := scenarioDB.NewRepository(dbPool)

	runtimeScenarioRepository := runtimescenario.NewScenarioRepository(dbPool)
	runtimeSessionRepository := runtimesession.NewSessionRepository(dbPool)
	runtimeStepRepository := runtimestep.NewStepRepository(dbPool)
	runtimeTokenRepository := runtimetoken.NewTestTokensRepository(dbPool)

	trackingConnection, err := dbPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire tracking database connection: %w", err)
	}
	defer trackingConnection.Release()

	trackingSessionRepository := trackingsession.NewSessionRepository(dbPool)
	trackingEventRepository := trackingevent.NewEventRepository(trackingConnection.Conn())
	trackingScenarioRepository := trackingscenario.NewScenarioRepository(dbPool)
	trackingStepRepository := trackingstep.NewStepRepository(dbPool)

	analyticsRepository := analyticsrepository.NewAnalyticsRepository(dbPool)

	keyGenerator := keygen.NewGenerator()
	testToken := testtoken.NewService()

	elementService := element.NewElementService(
		elementRepository,
		projectRepository,
		stepRepository,
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

	testTokenService := test_token.NewTestTokenService(
		testTokenRepository,
		scenarioRepository,
		testToken,
		time.Minute*30,
		logger,
	)

	stepService := step.NewStepService(
		stepRepository,
		scenarioRepository,
		elementRepository,
		txManager,
		logger,
	)

	scenarioService := scenario.NewScenarioService(
		scenarioRepository,
		stepRepository,
		projectRepository,
		txManager,
		logger,
	)

	runtimeService := runtimeservice.NewRuntimeService(
		runtimeScenarioRepository,
		runtimeSessionRepository,
		runtimeStepRepository,
		runtimeTokenRepository,
		testToken,
	)

	trackingService := trackingservice.NewTrackingService(
		trackingSessionRepository,
		trackingEventRepository,
		trackingScenarioRepository,
		trackingStepRepository,
	)

	analyticsService := analyticsservice.NewAnalyticsService(analyticsRepository)

	projectHandler := projecthttp.NewHandler(
		elementService,
		projectService,
		logger,
	)

	scenarioHandler := scenariohttp.NewHandler(
		scenarioService,
		stepService,
		testTokenService,
		logger,
	)

	runtimeHandler := runtimehttp.NewHandler(runtimeService)
	trackingHandler := trackinghttp.NewTrackingHandler(trackingService)
	analyticsHandler := analyticshttp.NewAnalyticsHandler(analyticsService)
	pdfHandler := pdfhttp.NewPDFHandler(analyticsService)

	server := httpserver.NewServer(
		cfg.HTTPConfig,
		platformmiddleware.CORS(cfg.HTTPConfig.AllowedOrigins),
	)

	server.RegisterRouteGroup(
		"/api/v1/",
		[]httpserver.Middleware{},
		projectHandler,
		scenarioHandler,
		analyticsHandler,
		pdfHandler,
	)

	server.RegisterRouteGroup(
		"/api/v1/sdk/",
		[]httpserver.Middleware{
			platformmiddleware.ExtractProjectKey,
			platformmiddleware.ExtractTestToken,
		},
		runtimeHandler,
		trackingHandler,
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
