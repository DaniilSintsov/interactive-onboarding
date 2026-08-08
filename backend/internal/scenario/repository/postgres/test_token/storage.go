package test_token

import (
	"context"
	"errors"
	"fmt"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/postgres/transactor"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/errs"
	sqlc "github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/repository/postgres/test_token/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	foreignKeyViolationCode = "23503"

	testTokenScenarioFKConstraint = "scenario_test_tokens_scenario_fk"
)

type testTokenRepository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *testTokenRepository {
	return &testTokenRepository{
		queries: sqlc.New(db),
	}
}

func (repo *testTokenRepository) getQueries(ctx context.Context) *sqlc.Queries {
	if tx, err := transactor.ExtractTx(ctx); err == nil {
		return repo.queries.WithTx(tx)
	}

	return repo.queries
}

func (repo *testTokenRepository) Create(
	ctx context.Context,
	token entity.ScenarioTestToken,
) (entity.ScenarioTestToken, error) {
	row, err := repo.getQueries(ctx).CreateScenarioTestToken(ctx, sqlc.CreateScenarioTestTokenParams{
		ScenarioID: token.ScenarioID,
		Hash:       token.Hash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  token.ExpiresAt,
			Valid: true,
		},
	},
	)
	if err != nil {
		if isConstraintError(err, foreignKeyViolationCode, testTokenScenarioFKConstraint) {
			return entity.ScenarioTestToken{}, errs.ErrScenarioNotFound
		}

		return entity.ScenarioTestToken{}, fmt.Errorf("scenario test token repository - create: %w", err)
	}

	return entity.ScenarioTestToken{
		ID:         row.ID,
		ScenarioID: row.ScenarioID,
		Hash:       row.Hash,
		CreatedAt:  row.CreatedAt.Time.UTC(),
		ExpiresAt:  row.ExpiresAt.Time.UTC(),
	}, nil
}

func isConstraintError(err error, code string, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) &&
		pgErr.Code == code &&
		pgErr.ConstraintName == constraint
}
