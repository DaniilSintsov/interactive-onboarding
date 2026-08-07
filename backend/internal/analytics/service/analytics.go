package service

import (
	"context"
	"time"

	"interactive-onboarding/internal/analytics/repository"
)

type AnalyticsService struct {
	repo *repository.AnalyticsRepository
}

func NewAnalyticsService(repo *repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

type ScenarioAnalytics struct {
	ScenarioID                   string  `json:"scenario_id"`
	Started                      int64   `json:"started"`
	Completed                    int64   `json:"completed"`
	Skipped                      int64   `json:"skipped"`
	CompletionRate               float64 `json:"completion_rate"`
	SkipRate                     float64 `json:"skip_rate"`
	AverageCompletionTimeSeconds float64 `json:"average_completion_time_seconds"`
}

type ProjectAnalytics struct {
	ProjectID         string  `json:"project_id"`
	TotalScenarios    int     `json:"total_scenarios"`
	EnabledScenarios  int     `json:"enabled_scenarios"`
	SessionsStarted   int64   `json:"sessions_started"`
	SessionsCompleted int64   `json:"sessions_completed"`
	SessionsSkipped   int64   `json:"sessions_skipped"`
	CompletionRate    float64 `json:"completion_rate"`
	SkipRate          float64 `json:"skip_rate"`
}

type StepAnalytics struct {
	StepID         string  `json:"step_id"`
	Position       int     `json:"position"`
	Title          string  `json:"title"`
	Shown          int64   `json:"shown"`
	Completed      int64   `json:"completed"`
	Skipped        int64   `json:"skipped"`
	CompletionRate float64 `json:"completion_rate"`
	SkipRate       float64 `json:"skip_rate"`
	DropOffRate    float64 `json:"drop_off_rate"`
}

type DetailedScenarioAnalytics struct {
	ScenarioAnalytics
	Steps []StepAnalytics `json:"steps"`
}

func (s *AnalyticsService) GetScenarioAnalytics(ctx context.Context, scenarioID string, from, to *time.Time) (ScenarioAnalytics, error) {
	stats, err := s.repo.GetScenarioStats(ctx, scenarioID, from, to)
	if err != nil {
		return ScenarioAnalytics{}, err
	}

	completionRate := 0.0
	if stats.Started > 0 {
		completionRate = float64(stats.Completed) / float64(stats.Started)
	}

	skipRate := 0.0
	if stats.Started > 0 {
		skipRate = float64(stats.Skipped) / float64(stats.Started)
	}

	return ScenarioAnalytics{
		ScenarioID:                   scenarioID,
		Started:                      stats.Started,
		Completed:                    stats.Completed,
		Skipped:                      stats.Skipped,
		CompletionRate:               completionRate,
		SkipRate:                     skipRate,
		AverageCompletionTimeSeconds: stats.AvgDurationSeconds,
	}, nil
}

func (s *AnalyticsService) GetDetailedScenarioAnalytics(ctx context.Context, scenarioID string, from, to *time.Time) (DetailedScenarioAnalytics, error) {
	summary, err := s.GetScenarioAnalytics(ctx, scenarioID, from, to)
	if err != nil {
		return DetailedScenarioAnalytics{}, err
	}

	stepStats, err := s.repo.GetStepFunnel(ctx, scenarioID, from, to)
	if err != nil {
		return DetailedScenarioAnalytics{}, err
	}

	steps := make([]StepAnalytics, 0, len(stepStats))
	for _, stat := range stepStats {
		completionRate := 0.0
		if stat.Shown > 0 {
			completionRate = float64(stat.Completed) / float64(stat.Shown)
		}
		skipRate := 0.0
		if stat.Shown > 0 {
			skipRate = float64(stat.Skipped) / float64(stat.Shown)
		}

		steps = append(steps, StepAnalytics{
			StepID:         stat.StepID,
			Position:       stat.Position,
			Title:          stat.Title,
			Shown:          stat.Shown,
			Completed:      stat.Completed,
			Skipped:        stat.Skipped,
			CompletionRate: completionRate,
			SkipRate:       skipRate,
			DropOffRate:    0,
		})
	}

	return DetailedScenarioAnalytics{
		ScenarioAnalytics: summary,
		Steps:             steps,
	}, nil
}

func (s *AnalyticsService) GetProjectAnalytics(ctx context.Context, projectID string, from, to *time.Time) (ProjectAnalytics, error) {
	total, enabled, started, completed, skipped, err := s.repo.GetProjectStats(ctx, projectID, from, to)
	if err != nil {
		return ProjectAnalytics{}, err
	}

	completionRate := 0.0
	if started > 0 {
		completionRate = float64(completed) / float64(started)
	}

	skipRate := 0.0
	if started > 0 {
		skipRate = float64(skipped) / float64(started)
	}

	return ProjectAnalytics{
		ProjectID:         projectID,
		TotalScenarios:    total,
		EnabledScenarios:  enabled,
		SessionsStarted:   started,
		SessionsCompleted: completed,
		SessionsSkipped:   skipped,
		CompletionRate:    completionRate,
		SkipRate:          skipRate,
	}, nil
}
