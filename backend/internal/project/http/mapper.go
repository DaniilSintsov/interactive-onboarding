package projecthttp

import (
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/port"
)

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
