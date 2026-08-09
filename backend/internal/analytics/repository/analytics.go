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
	DropOff   int64
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

func (r *AnalyticsRepository) ScenarioExistsActiveWithProject(ctx context.Context, scenarioID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 
			FROM onboarding.scenarios s
			JOIN onboarding.projects p ON p.id = s.project_id
			WHERE s.id = $1 
			AND s.deleted_at IS NULL
			AND p.deleted_at IS NULL
		)
	`, scenarioID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check scenario exists with active project: %w", err)
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
			COALESCE(AVG(GREATEST(0, EXTRACT(EPOCH FROM (finished_at - started_at)))) FILTER (WHERE status = 'completed'), 0) AS avg_duration_seconds
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
		),
		step_shown_sessions AS (
			SELECT 
				e.step_id,
				ARRAY_AGG(DISTINCT e.session_id) AS sessions
			FROM onboarding.events e
			WHERE e.type = 'step_shown'
			AND e.session_id IN (SELECT id FROM cohort_sessions)
			GROUP BY e.step_id
		),
		step_next_sessions AS (
			SELECT 
				s.id AS step_id,
				LEAD(s.id) OVER (ORDER BY s.step_num) AS next_step_id,
				LEAD(ss.sessions) OVER (ORDER BY s.step_num) AS next_sessions
			FROM onboarding.steps s
			LEFT JOIN step_shown_sessions ss ON ss.step_id = s.id
			WHERE s.scenario_id = $1 AND s.deleted_at IS NULL
		)
		SELECT 
			s.id AS step_id,
			s.step_num,
			s.title,
			COALESCE(CARDINALITY(ss.sessions), 0) AS shown,
			COALESCE(
				(SELECT COUNT(DISTINCT session_id) 
				 FROM onboarding.events 
				 WHERE step_id = s.id 
				 AND type = 'step_completed'
				 AND session_id IN (SELECT id FROM cohort_sessions)
				), 0
			) AS completed,
			COALESCE(
				(SELECT COUNT(DISTINCT session_id) 
				 FROM onboarding.events 
				 WHERE step_id = s.id 
				 AND type = 'step_skipped'
				 AND session_id IN (SELECT id FROM cohort_sessions)
				), 0
			) AS skipped,
			COALESCE(
				CASE 
					WHEN sns.next_step_id IS NULL THEN
						(SELECT COUNT(DISTINCT session_id)
						 FROM onboarding.events 
						 WHERE step_id = s.id 
						 AND type = 'step_shown'
						 AND session_id IN (SELECT id FROM cohort_sessions)
						 AND session_id NOT IN (
							SELECT id FROM onboarding.sessions 
							WHERE scenario_id = $1 
							AND status = 'completed'
							AND id IN (SELECT id FROM cohort_sessions)
						 )
						)
					ELSE
						CARDINALITY(
							ARRAY(
								SELECT UNNEST(ss.sessions)
								EXCEPT
								SELECT UNNEST(COALESCE(sns.next_sessions, ARRAY[]::uuid[]))
							)
						)
				END, 0
			) AS drop_off
		FROM onboarding.steps s
		LEFT JOIN step_shown_sessions ss ON ss.step_id = s.id
		LEFT JOIN step_next_sessions sns ON sns.step_id = s.id
		WHERE s.scenario_id = $1 AND s.deleted_at IS NULL
		GROUP BY s.id, s.step_num, s.title, ss.sessions, sns.next_step_id, sns.next_sessions
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
			&stat.DropOff,
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
			WHERE project_id = $1 AND deleted_at IS NULL
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
