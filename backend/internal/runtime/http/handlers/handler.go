package handlers

import (
	"net/http"

	runtimeModel "github.com/DaniilSintsov/interactive-onboarding/backend/internal/runtime/model"
	"github.com/go-playground/validator/v10"
)

type RuntimeService interface {
	FindScenario(pageId string, userId string) (*runtimeModel.RuntimeScenario, error)
}

type RuntimeHandler struct {
	service RuntimeService
}

var requestValidator = validator.New()

func NewHandler(srvc RuntimeService) *RuntimeHandler {
	return &RuntimeHandler{
		service: srvc,
	}
}

func (h *RuntimeHandler) GetScenario(w http.ResponseWriter, r *http.Request) {
	scenarioRequest := new(runtimeModel.ResolveScenarioRequest)
	if err := ParseJson(r, scenarioRequest); err != nil {
		raiseError(w, "failed to parse scenario request", err, http.StatusBadRequest)
		return
	}

	if err := requestValidator.Struct(scenarioRequest); err != nil {
		raiseError(w, "invalid scenario request", err, http.StatusUnprocessableEntity)
		return
	}
	scenarioResponse, err := h.service.FindScenario(scenarioRequest.Page, scenarioRequest.UserID)
	if err != nil {
		raiseError(w, "invalid scenario response", err, http.StatusBadRequest)
		return
	}
	if scenarioResponse == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_ = WriteJson(w, http.StatusOK, scenarioResponse)
}
