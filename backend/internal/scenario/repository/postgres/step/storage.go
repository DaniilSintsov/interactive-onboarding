package step

import (
	"context"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/postgres/transactor"
	sqlc "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/repository/postgres/element/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type stepRepository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *stepRepository {
	return &stepRepository{
		queries: sqlc.New(db),
	}
}

func (e *stepRepository) getQueries(ctx context.Context) *sqlc.Queries {
	if tx, err := transactor.ExtractTx(ctx); err == nil {
		return e.queries.WithTx(tx)
	}

	return e.queries
}

func textFromPtr(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{
			Valid: false,
		}
	}

	return pgtype.Text{
		String: *value,
		Valid:  true,
	}
}
