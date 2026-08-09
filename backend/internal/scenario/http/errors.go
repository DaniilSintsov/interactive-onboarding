package scenariohttp

import (
	"errors"
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
	projecterrs "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	scenarioerrs "github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/errs"
	"go.uber.org/zap"
)

func mapError(err error) httpserver.ErrorMapping {
	if mapping, ok := httpserver.MapRequestError(err); ok {
		return mapping
	}

	switch {
	case errors.Is(err, scenarioerrs.ErrProjectNotFound),
		errors.Is(err, projecterrs.ErrProjectNotFound):
		return httpserver.NewErrorMapping(http.StatusNotFound, "project_not_found", "project not found", false)
	case errors.Is(err, projecterrs.ErrElementNotFound):
		return httpserver.NewErrorMapping(http.StatusNotFound, "element_not_found", "element not found", false)
	case errors.Is(err, scenarioerrs.ErrScenarioNotFound):
		return httpserver.NewErrorMapping(http.StatusNotFound, "scenario_not_found", "scenario not found", false)
	case errors.Is(err, scenarioerrs.ErrStepNotFound):
		return httpserver.NewErrorMapping(http.StatusNotFound, "step_not_found", "step not found", false)
	case errors.Is(err, scenarioerrs.ErrScenarioImmutable):
		return httpserver.NewErrorMapping(http.StatusConflict, "scenario_immutable", "scenario cannot be modified in current status", false)
	case errors.Is(err, scenarioerrs.ErrInvalidScenarioStatusTransition):
		return httpserver.NewErrorMapping(http.StatusConflict, "invalid_status_transition", "invalid scenario status transition", false)
	case errors.Is(err, scenarioerrs.ErrStepNumberAlreadyExists):
		return httpserver.NewErrorMapping(http.StatusConflict, "step_number_already_exists", "step number already exists", false)
	case errors.Is(err, scenarioerrs.ErrStepDoesNotBelongToScenario):
		return httpserver.NewValidationErrorMapping("invalid_step_order", "step does not belong to scenario", "ordered_step_ids")
	case errors.Is(err, scenarioerrs.ErrScenarioProjectIDRequired):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_project_id", "project ID is required", false)
	case errors.Is(err, scenarioerrs.ErrStepScenarioIDRequired),
		errors.Is(err, scenarioerrs.ErrScenarioIDRequired):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_scenario_id", "scenario ID is required", false)
	case errors.Is(err, scenarioerrs.ErrScenarioUpdateFieldsRequired):
		return httpserver.NewValidationErrorMapping("scenario_update_fields_required", "at least one scenario field is required", "body")
	case errors.Is(err, scenarioerrs.ErrScenarioLimitInvalid):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_limit", "limit must be between 1 and 100", false)
	case errors.Is(err, scenarioerrs.ErrScenarioOffsetInvalid):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_offset", "offset must be a non-negative integer", false)
	case errors.Is(err, scenarioerrs.ErrStepElementIDRequired):
		return httpserver.NewValidationErrorMapping("step_element_id_required", "step element ID is required", "element_id")
	case errors.Is(err, scenarioerrs.ErrStepIDRequired):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_step_id", "step ID is required", false)
	case errors.Is(err, scenarioerrs.ErrStepTitleRequired):
		return httpserver.NewValidationErrorMapping("step_title_required", "step title is required", "title")
	case errors.Is(err, scenarioerrs.ErrInvalidStepNumber):
		return httpserver.NewValidationErrorMapping("invalid_step_number", "step number must be between 1 and 2147483647", "step_num")
	case errors.Is(err, scenarioerrs.ErrStepTitleTooLong):
		return httpserver.NewValidationErrorMapping("step_title_too_long", "step title is too long", "title")
	case errors.Is(err, scenarioerrs.ErrStepDescriptionTooLong):
		return httpserver.NewValidationErrorMapping("step_description_too_long", "step description is too long", "description")
	case errors.Is(err, scenarioerrs.ErrStepUpdateFieldsRequired):
		return httpserver.NewValidationErrorMapping("step_update_fields_required", "at least one step field is required", "body")
	case errors.Is(err, scenarioerrs.ErrStepOrderRequired):
		return httpserver.NewValidationErrorMapping("step_order_required", "ordered_step_ids must not be empty", "ordered_step_ids")
	case errors.Is(err, scenarioerrs.ErrScenarioNameRequired):
		return httpserver.NewValidationErrorMapping("scenario_name_required", "scenario name is required", "name")
	case errors.Is(err, scenarioerrs.ErrScenarioNameInvalid):
		return httpserver.NewValidationErrorMapping("scenario_name_invalid", "scenario name must not be empty", "name")
	case errors.Is(err, scenarioerrs.ErrScenarioPagePatternRequired):
		return httpserver.NewValidationErrorMapping("scenario_page_pattern_required", "scenario page pattern is required", "page_pattern")
	case errors.Is(err, scenarioerrs.ErrScenarioPagePatternInvalid):
		return httpserver.NewValidationErrorMapping("scenario_page_pattern_invalid", "scenario page pattern must not be empty", "page_pattern")
	case errors.Is(err, scenarioerrs.ErrScenarioNameTooLong):
		return httpserver.NewValidationErrorMapping("scenario_name_too_long", "scenario name is too long", "name")
	case errors.Is(err, scenarioerrs.ErrScenarioDescriptionTooLong):
		return httpserver.NewValidationErrorMapping("scenario_description_too_long", "scenario description is too long", "description")
	case errors.Is(err, scenarioerrs.ErrScenarioPagePatternTooLong):
		return httpserver.NewValidationErrorMapping("scenario_page_pattern_too_long", "scenario page pattern is too long", "page_pattern")
	case errors.Is(err, scenarioerrs.ErrScenarioStatusUnknown):
		return httpserver.NewValidationErrorMapping("unknown_scenario_status", "unknown scenario status", "status")
	case errors.Is(err, scenarioerrs.ErrScenarioTestTokenScenarioIDRequired):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_scenario_id", "scenario ID is required", false)
	case errors.Is(err, scenarioerrs.ErrScenarioTestTokenHashInvalid):
		return httpserver.NewValidationErrorMapping("invalid_test_token_hash", "test token hash is invalid", "token")
	case errors.Is(err, scenarioerrs.ErrScenarioTestTokenExpirationInvalid):
		return httpserver.NewValidationErrorMapping("invalid_test_token_expiration", "test token expiration is invalid", "expires_at")
	default:
		return httpserver.InternalErrorMapping()
	}
}

func (h *Handler) handleError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	mapping := mapError(err)
	if mapping.Log {
		h.logger.Error(
			"scenario http handler failed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)
	}

	httpserver.WriteJSON(w, mapping.Status, mapping.Response)
}
