package projecthttp

import (
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/port"
	"github.com/google/uuid"
)

func createProjectRequestToParams(request CreateProjectRequest) port.CreateProjectParams {
	elements := make([]port.CreateProjectElementParams, 0, len(request.Elements))
	for _, element := range request.Elements {
		elements = append(elements, port.CreateProjectElementParams{
			Key:         element.Key,
			Label:       element.Label,
			Description: element.Description,
		})
	}

	return port.CreateProjectParams{
		Name:     request.Name,
		Elements: elements,
	}
}

func createElementRequestToParams(request CreateElementRequest, projectID uuid.UUID) port.CreateElementParams {
	return port.CreateElementParams{
		ProjectID:   projectID,
		Key:         request.Key,
		Label:       request.Label,
		Description: request.Description,
	}
}

func updateElementRequestToParams(request UpdateElementRequest, projectID uuid.UUID, elementID uuid.UUID) port.UpdateElementParams {
	return port.UpdateElementParams{
		ProjectID:   projectID,
		ElementID:   elementID,
		Key:         request.Key,
		Label:       request.Label,
		Description: request.Description,
	}
}

func projectToResponse(project entity.Project) ProjectResponse {
	return ProjectResponse{
		ID:         project.ID,
		Name:       project.Name,
		ProjectKey: project.ProjectKey,
		CreatedAt:  project.CreatedAt,
		UpdatedAt:  project.UpdatedAt,
	}
}

func elementToResponse(element entity.Element) ElementResponse {
	return ElementResponse{
		ID:          element.ID,
		ProjectID:   element.ProjectID,
		Key:         element.Key,
		Label:       element.Label,
		Description: element.Description,
		CreatedAt:   element.CreatedAt,
		UpdatedAt:   element.UpdatedAt,
	}
}

func projectWithElementsToResponse(
	result port.ProjectWithElements,
) ProjectWithElementsResponse {
	elements := make([]ElementResponse, 0, len(result.Elements))

	for _, element := range result.Elements {
		elements = append(elements, elementToResponse(element))
	}

	return ProjectWithElementsResponse{
		ProjectResponse: projectToResponse(result.Project),
		Elements:        elements,
	}
}

func projectListToResponse(result port.ListProjectsResult) ProjectListResponse {
	items := make([]ProjectResponse, 0, len(result.Projects))
	for _, project := range result.Projects {
		items = append(items, projectToResponse(project))
	}

	return ProjectListResponse{
		Items:  items,
		Total:  result.Total,
		Limit:  result.Limit,
		Offset: result.Offset,
	}
}

func elementsToResponse(elements []entity.Element) []ElementResponse {
	response := make([]ElementResponse, 0, len(elements))

	for _, element := range elements {
		response = append(response, elementToResponse(element))
	}

	return response
}
