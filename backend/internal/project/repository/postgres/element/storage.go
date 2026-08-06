package element

import (
	"context"
	"errors"
	"fmt"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	sqlc "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/repository/postgres/element/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type elementRepository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *elementRepository {
	return &elementRepository{
		queries: sqlc.New(db),
	}
}

const (
	uniqueViolationCode     = "23505"
	foreignKeyViolationCode = "23503"

	elementKeyUniqueConstraint = "elements_project_id_key_unique"
	elementProjectFKConstraint = "elements_project_fk"
)

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

func (repo *elementRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]entity.Element, error) {
	rows, err := repo.queries.ListElementsByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("element repository - list by project id: %w", err)
	}

	elements := make([]entity.Element, 0, len(rows))
	for _, row := range rows {
		elements = append(elements, entity.Element{
			ID:          row.ID,
			ProjectID:   row.ProjectID,
			Key:         row.Key,
			Label:       row.Label,
			Description: row.Description,
			CreatedAt:   row.CreatedAt.Time.UTC(),
			UpdatedAt:   row.UpdatedAt.Time.UTC(),
		})
	}

	return elements, nil
}

func (repo *elementRepository) Create(
	ctx context.Context,
	element entity.Element,
) (entity.Element, error) {
	createdElement, err := repo.queries.CreateElement(ctx, sqlc.CreateElementParams{
		ProjectID:   element.ProjectID,
		Key:         element.Key,
		Label:       element.Label,
		Description: element.Description,
	})
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == uniqueViolationCode &&
				pgErr.ConstraintName == elementKeyUniqueConstraint:

				return entity.Element{}, errs.ErrElementKeyAlreadyExists

			case pgErr.Code == foreignKeyViolationCode &&
				pgErr.ConstraintName == elementProjectFKConstraint:

				return entity.Element{}, errs.ErrProjectNotFound
			}
		}

		return entity.Element{}, fmt.Errorf("element repository - create: %w", err)
	}

	return entity.Element{
		ID:          createdElement.ID,
		ProjectID:   createdElement.ProjectID,
		Key:         createdElement.Key,
		Label:       createdElement.Label,
		Description: createdElement.Description,
		CreatedAt:   createdElement.CreatedAt.Time.UTC(),
		UpdatedAt:   createdElement.UpdatedAt.Time.UTC(),
	}, nil
}

type UpdateElementParams struct {
	ProjectID   uuid.UUID
	ElementID   uuid.UUID
	Key         *string
	Label       *string
	Description *string
}

// Update TODO: использовать параметры чтобы сделать возможным передавать часть полей
func (repo *elementRepository) Update(
	ctx context.Context,
	element entity.Element,
) (entity.Element, error) {
	updatedElement, err := repo.queries.UpdateElement(ctx, sqlc.UpdateElementParams{
		ProjectID:   element.ProjectID,
		ElementID:   element.ID,
		Key:         textFromPtr(&element.Key),
		Label:       textFromPtr(&element.Label),
		Description: textFromPtr(&element.Description),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Element{}, errs.ErrElementNotFound
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == uniqueViolationCode &&
			pgErr.ConstraintName == elementKeyUniqueConstraint {

			return entity.Element{}, errs.ErrElementKeyAlreadyExists
		}

		return entity.Element{}, fmt.Errorf("element repository - update: %w", err)
	}

	return entity.Element{
		ID:          updatedElement.ID,
		ProjectID:   updatedElement.ProjectID,
		Key:         updatedElement.Key,
		Label:       updatedElement.Label,
		Description: updatedElement.Description,
		CreatedAt:   updatedElement.CreatedAt.Time.UTC(),
		UpdatedAt:   updatedElement.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *elementRepository) Delete(
	ctx context.Context,
	projectID uuid.UUID,
	elementID uuid.UUID,
) error {
	_, err := repo.queries.DeleteElement(ctx, sqlc.DeleteElementParams{
		ProjectID: projectID,
		ElementID: elementID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrElementNotFound
		}

		return fmt.Errorf("element repository - delete: %w", err)
	}

	return nil
}
