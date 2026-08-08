package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	runtimeModel "github.com/DaniilSintsov/interactive-onboarding/backend/internal/runtime/model"
)

var (
	errTestTokenIsEmpty       = errors.New("test token is empty")
	ErrScenarioNotFound       = errors.New("scenario was not found")
	ErrProjectMismatch        = errors.New("scenario belongs to another project")
	ErrTokenIsExpired         = errors.New("token is expired")
	ErrProjectTokenIsNotValid = errors.New("project token is not valid")
)

type (
	ScenarioRepository interface {
		GetScenarioById(ctx context.Context, scenarioId string) (*runtimeModel.Scenario, error)
		GetScenarioByIdAndProjectKey(ctx context.Context, scenarioId, projectKey string) (*runtimeModel.Scenario, error)
		GetScenariosByPagePatternAndProjectkey(ctx context.Context, pagePattern, projectKey string) ([]runtimeModel.Scenario, error)
	}

	SessionRepository interface {
		GetSessionByScenarioAndUser(ctx context.Context, scenarioId, userId string) (*runtimeModel.Session, error)
	}

	StepRepository interface {
		GetStepsByScenarioId(ctx context.Context, scenarioId string) (*runtimeModel.RuntimeScenario, error)
	}
	TestTokensRepository interface {
		GetTokenByHash(ctx context.Context, hash []byte) (*runtimeModel.TestToken, error)
	}
	TokenHasher interface {
		Hash(rawToken string) []byte
	}
)

type RuntimeService struct {
	scenario    ScenarioRepository
	session     SessionRepository
	steps       StepRepository
	tokens      TestTokensRepository
	tokenHasher TokenHasher
}

func NewRuntimeService(
	rep ScenarioRepository, ses SessionRepository, st StepRepository, tok TestTokensRepository, hasher TokenHasher,
) *RuntimeService {
	return &RuntimeService{
		scenario:    rep,
		session:     ses,
		steps:       st,
		tokens:      tok,
		tokenHasher: hasher,
	}
}

func (r *RuntimeService) FindScenarios(ctx context.Context, pagePattern, userId string) (*runtimeModel.RuntimeScenarioResolveResponse, error) {
	response := new(runtimeModel.RuntimeScenarioResolveResponse)
	testScenario, err := r.checkTestToken(ctx)
	if err != nil && !errors.Is(err, errTestTokenIsEmpty) {
		return nil, err
	} else if err == nil {
		response.IsTest = true
		response.Scenarios = []runtimeModel.RuntimeScenario{*testScenario}
		return response, nil
	}

	projectKey, ok := ctx.Value("projectKey").(string)
	if !ok {
		return nil, ErrProjectTokenIsNotValid
	}
	scenarios, err := r.scenario.GetScenariosByPagePatternAndProjectkey(ctx, pagePattern, projectKey)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(scenarios); i++ {
		session, err := r.session.GetSessionByScenarioAndUser(ctx, scenarios[i].ID, userId)
		if err == nil {
			if session.Status != runtimeModel.SessionStatusActive {
				sixMonthsAgo := time.Now().AddDate(0, -6, 0)
				if session.FinishedAt.Before(sixMonthsAgo) {
					continue
				}
				scenarios[i] = scenarios[len(scenarios)-1]
				scenarios = scenarios[:len(scenarios)-1]
				i--
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	for i := range scenarios {
		runtimeScenario, err := r.steps.GetStepsByScenarioId(ctx, scenarios[i].ID)
		if err != nil {
			return nil, err
		}
		response.Scenarios = append(response.Scenarios, *runtimeScenario)
	}

	return response, nil
}

func (r *RuntimeService) checkTestToken(ctx context.Context) (*runtimeModel.RuntimeScenario, error) {
	testToken, ok := ctx.Value("testToken").(string)
	if !ok || testToken == "" {
		return nil, errTestTokenIsEmpty
	}
	hash := r.tokenHasher.Hash(testToken)
	token, err := r.tokens.GetTokenByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrScenarioNotFound
		}
		return nil, err
	}
	if token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenIsExpired
	}
	err = r.validateProjectKey(ctx, token.ScenarioID)
	if err != nil {
		return nil, ErrProjectMismatch
	}
	return r.steps.GetStepsByScenarioId(ctx, token.ScenarioID)
}

func (r *RuntimeService) validateProjectKey(ctx context.Context, scenarioId string) error {
	projectKey, ok := ctx.Value("projectKey").(string)
	if !ok {
		return ErrProjectTokenIsNotValid
	}
	_, err := r.scenario.GetScenarioByIdAndProjectKey(ctx, scenarioId, projectKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrScenarioNotFound
		} else {
			return err
		}
	}
	return nil
}
