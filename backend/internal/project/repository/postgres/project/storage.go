package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/postgres/transactor"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	sqlc "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/repository/postgres/project/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolationCode = "23505"

	projectKeyUniqueConstraint = "projects_project_key_unique"
)

type projectRepository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *projectRepository {
	return &projectRepository{
		queries: sqlc.New(db),
	}
}

func (p *projectRepository) getQueries(ctx context.Context) *sqlc.Queries {
	if tx, err := transactor.ExtractTx(ctx); err == nil {
		return p.queries.WithTx(tx)
	}

	return p.queries
}

func (repo *projectRepository) Create(
	ctx context.Context,
	project entity.Project,
) (entity.Project, error) {
	createdProject, err := repo.getQueries(ctx).CreateProject(ctx, sqlc.CreateProjectParams{
		Name:       project.Name,
		ProjectKey: project.ProjectKey,
	})
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == uniqueViolationCode &&
			pgErr.ConstraintName == projectKeyUniqueConstraint {

			return entity.Project{}, errs.ErrProjectKeyAlreadyExists
		}

		return entity.Project{}, fmt.Errorf("project repository - create: %w", err)
	}

	return entity.Project{
		ID:         createdProject.ID,
		Name:       createdProject.Name,
		ProjectKey: createdProject.ProjectKey,
		CreatedAt:  createdProject.CreatedAt.Time.UTC(),
		UpdatedAt:  createdProject.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *projectRepository) Update(
	ctx context.Context,
	projectID uuid.UUID,
	name string,
) (entity.Project, error) {
	updatedProject, err := repo.getQueries(ctx).UpdateProject(ctx, sqlc.UpdateProjectParams{
		Name:      name,
		ProjectID: projectID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Project{}, errs.ErrProjectNotFound
		}

		return entity.Project{}, fmt.Errorf("project repository - update: %w", err)
	}

	return entity.Project{
		ID:         updatedProject.ID,
		Name:       updatedProject.Name,
		ProjectKey: updatedProject.ProjectKey,
		CreatedAt:  updatedProject.CreatedAt.Time.UTC(),
		UpdatedAt:  updatedProject.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *projectRepository) Delete(
	ctx context.Context,
	projectID uuid.UUID,
) error {
	_, err := repo.getQueries(ctx).DeleteProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrProjectNotFound
		}

		return fmt.Errorf("project repository - delete: %w", err)
	}

	return nil
}

func (repo *projectRepository) List(
	ctx context.Context,
	limit int32,
	offset int32,
) ([]entity.Project, int64, error) {
	rows, err := repo.getQueries(ctx).ListProjects(ctx, sqlc.ListProjectsParams{
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("project repository - list: %w", err)
	}

	projects := make([]entity.Project, 0, len(rows))
	var total int64

	for _, row := range rows {
		total = row.Total
		if !row.ID.Valid {
			continue
		}

		projects = append(projects, entity.Project{
			ID:         uuid.UUID(row.ID.Bytes),
			Name:       row.Name.String,
			ProjectKey: row.ProjectKey.String,
			CreatedAt:  row.CreatedAt.Time.UTC(),
			UpdatedAt:  row.UpdatedAt.Time.UTC(),
		})
	}

	return projects, total, nil
}

func (repo *projectRepository) GetByID(
	ctx context.Context,
	projectID uuid.UUID,
) (entity.Project, error) {
	project, err := repo.getQueries(ctx).GetProjectByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Project{}, errs.ErrProjectNotFound
		}

		return entity.Project{}, fmt.Errorf("project repository - get by id: %w", err)
	}

	return entity.Project{
		ID:         project.ID,
		Name:       project.Name,
		ProjectKey: project.ProjectKey,
		CreatedAt:  project.CreatedAt.Time.UTC(),
		UpdatedAt:  project.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *projectRepository) GetIDByKey(
	ctx context.Context,
	projectKey string,
) (uuid.UUID, error) {
	id, err := repo.getQueries(ctx).GetProjectIDByKey(ctx, projectKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errs.ErrProjectNotFound
		}

		return uuid.Nil, fmt.Errorf("project repository - get id by key: %w", err)
	}

	return id, nil
}

func (repo *projectRepository) LockActive(
	ctx context.Context,
	projectID uuid.UUID,
) error {
	_, err := repo.getQueries(ctx).LockActiveProject(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrProjectNotFound
		}

		return fmt.Errorf("project repository - lock active: %w", err)
	}
	return nil
}
