package test_token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/errs"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	testTokenRepository interface {
		Create(ctx context.Context, token entity.ScenarioTestToken) (entity.ScenarioTestToken, error)
	}

	testTokenGenerator interface {
		Generate() (rawToken string, hash []byte, err error)
	}

	scenarioChecker interface {
		GetByID(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
	}
)

type testTokenService struct {
	testTokenRepository testTokenRepository
	scenarioChecker     scenarioChecker
	tokenGenerator      testTokenGenerator
	tokenTTL            time.Duration
	logger              *zap.Logger
}

func NewTestTokenService(
	testTokenRepository testTokenRepository,
	scenarioChecker scenarioChecker,
	tokenGenerator testTokenGenerator,
	tokenTTL time.Duration,
	logger *zap.Logger,
) *testTokenService {
	return &testTokenService{
		testTokenRepository: testTokenRepository,
		scenarioChecker:     scenarioChecker,
		tokenGenerator:      tokenGenerator,
		tokenTTL:            tokenTTL,
		logger:              logger,
	}
}

func (service *testTokenService) Create(
	ctx context.Context,
	scenarioID uuid.UUID,
) (port.CreatedScenarioTestToken, error) {
	_, err := service.scenarioChecker.GetByID(ctx, scenarioID)
	if err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) {
			return port.CreatedScenarioTestToken{}, err
		}
		return port.CreatedScenarioTestToken{}, service.wrapCreateError(err, scenarioID)
	}

	token, hash, err := service.tokenGenerator.Generate()
	if err != nil {
		return port.CreatedScenarioTestToken{}, service.wrapCreateError(err, scenarioID)
	}

	now := time.Now().UTC()
	testToken := entity.ScenarioTestToken{
		ScenarioID: scenarioID,
		Hash:       hash,
		CreatedAt:  now,
		ExpiresAt:  now.Add(service.tokenTTL),
	}

	if err = testToken.Validate(); err != nil {
		return port.CreatedScenarioTestToken{}, err
	}

	createdToken, err := service.testTokenRepository.Create(ctx, testToken)
	if err != nil {
		return port.CreatedScenarioTestToken{}, service.wrapCreateError(err, scenarioID)
	}

	return port.CreatedScenarioTestToken{
		RawToken:  token,
		ExpiresAt: createdToken.ExpiresAt,
	}, nil
}

func (service *testTokenService) wrapCreateError(err error, scenarioID uuid.UUID) error {
	service.logger.Error("scenario test token usecase - create failed",
		zap.String("scenario_id", scenarioID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("scenario test token usecase - create: scenario_id=%v: %w", scenarioID, err)
}
