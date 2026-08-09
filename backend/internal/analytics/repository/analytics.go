package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepository struct {
	db *pgxpool.Pool
}

func NewAnalyticsRepository(db *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

type ScenarioStats struct {
	Started            int64
	Completed          int64
	Skipped            int64
	AvgDurationSeconds float64
}

type StepStats struct {
	StepID    string
	Position  int
	Title     string
	Shown     int64
	Completed int64
	Skipped   int64
}

func (r *AnalyticsRepository) ScenarioExistsPhysical(ctx context.Context, scenarioID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM onboarding.scenarios WHERE id = $1)
	`, scenarioID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check scenario exists: %w", err)
	}
	return exists, nil
}

func (r *AnalyticsRepository) ScenarioExistsActive(ctx context.Context, scenarioID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM onboarding.scenarios WHERE id = $1 AND deleted_at IS NULL)
	`, scenarioID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check scenario exists active: %w", err)
	}
	return exists, nil
}

func (r *AnalyticsRepository) ProjectExistsPhysical(ctx context.Context, projectID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM onboarding.projects WHERE id = $1)
	`, projectID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check project exists: %w", err)
	}
	return exists, nil
}

func (r *AnalyticsRepository) ProjectExistsActive(ctx context.Context, projectID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM onboarding.projects WHERE id = $1 AND deleted_at IS NULL)
	`, projectID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check project exists active: %w", err)
	}
	return exists, nil
}

func (r *AnalyticsRepository) GetScenarioStats(ctx context.Context, scenarioID string, from, to *time.Time) (ScenarioStats, error) {
	query := `
		SELECT 
			COUNT(*) AS started,
			COUNT(*) FILTER (WHERE status = 'completed') AS completed,
			COUNT(*) FILTER (WHERE status = 'skipped') AS skipped,
			COALESCE(AVG(EXTRACT(EPOCH FROM (finished_at - started_at))) FILTER (WHERE status = 'completed'), 0) AS avg_duration_seconds
		FROM onboarding.sessions
		WHERE scenario_id = $1
		AND ($2::TIMESTAMPTZ IS NULL OR started_at >= $2)
		AND ($3::TIMESTAMPTZ IS NULL OR started_at < $3)
	`

	var stats ScenarioStats
	err := r.db.QueryRow(ctx, query, scenarioID, from, to).Scan(
		&stats.Started,
		&stats.Completed,
		&stats.Skipped,
		&stats.AvgDurationSeconds,
	)
	if err != nil {
		return ScenarioStats{}, fmt.Errorf("get scenario stats: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetStepFunnel(ctx context.Context, scenarioID string, from, to *time.Time) ([]StepStats, error) {
	query := `
		WITH cohort_sessions AS (
			SELECT id FROM onboarding.sessions
			WHERE scenario_id = $1
			AND ($2::TIMESTAMPTZ IS NULL OR started_at >= $2)
			AND ($3::TIMESTAMPTZ IS NULL OR started_at < $3)
		)
		SELECT 
			s.id AS step_id,
			s.step_num,
			s.title,
			COUNT(DISTINCT CASE WHEN e.type = 'step_shown' THEN e.session_id END) AS shown,
			COUNT(DISTINCT CASE WHEN e.type = 'step_completed' THEN e.session_id END) AS completed,
			COUNT(DISTINCT CASE WHEN e.type = 'step_skipped' THEN e.session_id END) AS skipped
		FROM onboarding.steps s
		LEFT JOIN onboarding.events e ON e.step_id = s.id
			AND e.session_id IN (SELECT id FROM cohort_sessions)
		WHERE s.scenario_id = $1
		GROUP BY s.id, s.step_num, s.title
		ORDER BY s.step_num ASC
	`

	rows, err := r.db.Query(ctx, query, scenarioID, from, to)
	if err != nil {
		return nil, fmt.Errorf("get step funnel: %w", err)
	}
	defer rows.Close()

	var steps []StepStats
	for rows.Next() {
		var stat StepStats
		if err := rows.Scan(
			&stat.StepID,
			&stat.Position,
			&stat.Title,
			&stat.Shown,
			&stat.Completed,
			&stat.Skipped,
		); err != nil {
			return nil, fmt.Errorf("scan step funnel: %w", err)
		}
		steps = append(steps, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return steps, nil
}

func (r *AnalyticsRepository) GetProjectStats(ctx context.Context, projectID string, from, to *time.Time) (totalScenarios int, enabledScenarios int, started int64, completed int64, skipped int64, err error) {
	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM onboarding.scenarios 
		WHERE project_id = $1 AND deleted_at IS NULL
	`, projectID).Scan(&totalScenarios)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("count total scenarios: %w", err)
	}

	err = r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM onboarding.scenarios 
		WHERE project_id = $1 AND status = 'enabled' AND deleted_at IS NULL
	`, projectID).Scan(&enabledScenarios)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("count enabled scenarios: %w", err)
	}

	query := `
		SELECT 
			COUNT(*) AS started,
			COUNT(*) FILTER (WHERE status = 'completed') AS completed,
			COUNT(*) FILTER (WHERE status = 'skipped') AS skipped
		FROM onboarding.sessions
		WHERE scenario_id IN (
			SELECT id FROM onboarding.scenarios 
			WHERE project_id = $1
		)
		AND ($2::TIMESTAMPTZ IS NULL OR started_at >= $2)
		AND ($3::TIMESTAMPTZ IS NULL OR started_at < $3)
	`

	err = r.db.QueryRow(ctx, query, projectID, from, to).Scan(&started, &completed, &skipped)
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("get project stats: %w", err)
	}
	return totalScenarios, enabledScenarios, started, completed, skipped, nil
}
