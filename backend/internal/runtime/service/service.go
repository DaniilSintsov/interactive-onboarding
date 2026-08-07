package service

import runtimeModel "github.com/DaniilSintsov/interactive-onboarding/backend/internal/runtime/model"

type ScenarioRepository interface {
	GetScenarioByPageId(pageId string) (runtimeModel.RuntimeScenario, error)
	GetUserById(userId string) (runtimeModel.User, error)
}

type RuntimeService struct {
	repository ScenarioRepository
}

func NewRuntimeService(rep ScenarioRepository) *RuntimeService {
	return &RuntimeService{
		repository: rep,
	}
}

func (r RuntimeService) FindScenario(pageId string, userId string) (*runtimeModel.RuntimeScenario, error) {
	scenario, err := r.repository.GetScenarioByPageId(pageId)
	if err != nil {
		return nil, err
	}

	user, err := r.repository.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user.Onboarded {
		// Maps to response code 204
		return nil, nil
	}

	return &scenario, nil
}
