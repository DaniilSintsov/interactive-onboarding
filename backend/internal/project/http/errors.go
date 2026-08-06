package projecthttp

import (
	"errors"
	"net/http"

	projecterrs "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	"go.uber.org/zap"
)

type errorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type errorMapping struct {
	Status   int
	Response errorResponse
	Log      bool
}

func mapError(err error) errorMapping {
	switch {
	case errors.Is(err, projecterrs.ErrProjectNotFound):
		return newErrorMapping(http.StatusNotFound, "project_not_found", "project not found", false)

	case errors.Is(err, projecterrs.ErrElementNotFound):
		return newErrorMapping(http.StatusNotFound, "element_not_found", "element not found", false)

	case errors.Is(err, projecterrs.ErrElementKeyAlreadyExists):
		return newErrorMapping(http.StatusConflict, "element_key_already_exists", "element key already exists", false)

	case errors.Is(err, projecterrs.ErrElementInUse):
		return newErrorMapping(http.StatusConflict, "element_in_use", "element is used by an active step", false)

	case errors.Is(err, projecterrs.ErrProjectIDRequired),
		errors.Is(err, projecterrs.ErrElementProjectIDRequired):
		return newErrorMapping(http.StatusBadRequest, "invalid_project_id", "project ID is required", false)

	case errors.Is(err, projecterrs.ErrElementIDRequired):
		return newErrorMapping(http.StatusBadRequest, "invalid_element_id", "element ID is required", false)

	case errors.Is(err, projecterrs.ErrLimitInvalid):
		return newErrorMapping(http.StatusBadRequest, "invalid_limit", "limit must be between 1 and 100", false)

	case errors.Is(err, projecterrs.ErrOffsetInvalid):
		return newErrorMapping(http.StatusBadRequest, "invalid_offset", "offset must be greater than or equal to 0", false)

	case errors.Is(err, projecterrs.ErrProjectNameRequired):
		return newValidationErrorMapping("project_name_required", "project name is required", "name")

	case errors.Is(err, projecterrs.ErrProjectNameTooLong):
		return newValidationErrorMapping("project_name_too_long", "project name is too long", "name")

	case errors.Is(err, projecterrs.ErrElementKeyRequired):
		return newValidationErrorMapping("element_key_required", "element key is required", "key")

	case errors.Is(err, projecterrs.ErrElementKeyTooLong):
		return newValidationErrorMapping("element_key_too_long", "element key is too long", "key")

	case errors.Is(err, projecterrs.ErrElementLabelRequired):
		return newValidationErrorMapping("element_label_required", "element label is required", "label")

	case errors.Is(err, projecterrs.ErrElementLabelTooLong):
		return newValidationErrorMapping("element_label_too_long", "element label is too long", "label")

	case errors.Is(err, projecterrs.ErrElementDescriptionTooLong):
		return newValidationErrorMapping("element_description_too_long", "element description is too long", "description")

	case errors.Is(err, projecterrs.ErrEmptyElementUpdateParams):
		return newErrorMapping(http.StatusUnprocessableEntity, "empty_update", "at least one field must be provided", false)

	case errors.Is(err, projecterrs.ErrFailedGenerateUniqueKey):
		return newErrorMapping(
			http.StatusInternalServerError,
			"internal_error",
			"internal server error",
			true,
		)

	default:
		return newErrorMapping(
			http.StatusInternalServerError,
			"internal_error",
			"internal server error",
			true,
		)
	}
}

func newErrorMapping(status int, code, message string, logError bool) errorMapping {
	return errorMapping{
		Status: status,
		Response: errorResponse{
			Code:    code,
			Message: message,
		},
		Log: logError,
	}
}

func newValidationErrorMapping(code, message, field string) errorMapping {
	mapping := newErrorMapping(http.StatusUnprocessableEntity, code, message, false)
	mapping.Response.Details = map[string]any{"field": field}

	return mapping
}

func (h *Handler) handleError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	mapping := mapError(err)

	if mapping.Log {
		h.logger.Error(
			"http handler failed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)
	}

	writeJSON(w, mapping.Status, mapping.Response)
}
