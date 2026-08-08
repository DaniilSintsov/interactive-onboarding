package scenariohttp

import (
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/port"
	"github.com/google/uuid"
)

func createScenarioRequestToParams(
	request CreateScenarioRequest,
	projectID uuid.UUID,
) port.CreateScenarioParams {
	return port.CreateScenarioParams{
		ProjectID:   projectID,
		Name:        request.Name,
		Description: request.Description,
		PagePattern: request.PagePattern,
	}
}

func updateScenarioRequestToParams(
	request UpdateScenarioRequest,
	scenarioID uuid.UUID,
) port.UpdateScenarioParams {
	return port.UpdateScenarioParams{
		ScenarioID:  scenarioID,
		Name:        request.Name,
		Description: request.Description,
		PagePattern: request.PagePattern,
	}
}

func createStepRequestToParams(
	request CreateStepRequest,
	scenarioID uuid.UUID,
) port.CreateStepParams {
	return port.CreateStepParams{
		ScenarioID:   scenarioID,
		ElementID:    request.ElementID,
		StepNum:      request.StepNum,
		Title:        request.Title,
		Description:  request.Description,
		FrontendData: request.FrontendData,
	}
}

func updateStepRequestToParams(
	request UpdateStepRequest,
	scenarioID uuid.UUID,
	stepID uuid.UUID,
) port.UpdateStepParams {
	return port.UpdateStepParams{
		ScenarioID:   scenarioID,
		StepID:       stepID,
		ElementID:    request.ElementID,
		Title:        request.Title,
		Description:  request.Description,
		FrontendData: request.FrontendData,
	}
}

func reorderStepsRequestToParams(
	request ReorderStepsRequest,
	scenarioID uuid.UUID,
) port.ReorderStepsParams {
	return port.ReorderStepsParams{
		ScenarioID:     scenarioID,
		OrderedStepIDs: request.OrderedStepIDs,
	}
}

func scenarioToResponse(scenario entity.Scenario) ScenarioResponse {
	return ScenarioResponse{
		ID:          scenario.ID,
		ProjectID:   scenario.ProjectID,
		Name:        scenario.Name,
		Description: scenario.Description,
		PagePattern: scenario.PagePattern,
		Status:      scenario.Status,
		PublishedAt: scenario.PublishedAt,
		CreatedAt:   scenario.CreatedAt,
		UpdatedAt:   scenario.UpdatedAt,
	}
}

func scenarioSummaryToResponse(summary port.ScenarioSummary) ScenarioSummaryResponse {
	return ScenarioSummaryResponse{
		ID:          summary.Scenario.ID,
		ProjectID:   summary.Scenario.ProjectID,
		Name:        summary.Scenario.Name,
		PagePattern: summary.Scenario.PagePattern,
		Status:      summary.Scenario.Status,
		StepsCount:  summary.StepsCount,
		PublishedAt: summary.Scenario.PublishedAt,
		CreatedAt:   summary.Scenario.CreatedAt,
		UpdatedAt:   summary.Scenario.UpdatedAt,
	}
}

func scenarioListToResponse(result port.ListScenariosResult) ScenarioListResponse {
	items := make([]ScenarioSummaryResponse, 0, len(result.Scenarios))
	for _, scenario := range result.Scenarios {
		items = append(items, scenarioSummaryToResponse(scenario))
	}

	return ScenarioListResponse{
		Items:  items,
		Total:  result.Total,
		Limit:  result.Limit,
		Offset: result.Offset,
	}
}

func stepToResponse(step entity.Step) StepResponse {
	return StepResponse{
		ID:           step.ID,
		ScenarioID:   step.ScenarioID,
		ElementID:    step.ElementID,
		StepNum:      step.StepNum,
		Title:        step.Title,
		Description:  step.Description,
		FrontendData: step.FrontendData,
		CreatedAt:    step.CreatedAt,
		UpdatedAt:    step.UpdatedAt,
	}
}

func stepsToResponse(steps []entity.Step) []StepResponse {
	response := make([]StepResponse, 0, len(steps))
	for _, step := range steps {
		response = append(response, stepToResponse(step))
	}
	return response
}

func scenarioWithStepsToResponse(result port.ScenarioWithSteps) ScenarioWithStepsResponse {
	return ScenarioWithStepsResponse{
		ScenarioResponse: scenarioToResponse(result.Scenario),
		Steps:            stepsToResponse(result.Steps),
	}
}

func testTokenToResponse(token port.CreatedScenarioTestToken) ScenarioTestTokenResponse {
	return ScenarioTestTokenResponse{
		Token:     token.RawToken,
		ExpiresAt: token.ExpiresAt,
	}
}
