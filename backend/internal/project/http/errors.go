package projecthttp

import (
	"errors"
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
	projecterrs "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	"go.uber.org/zap"
)

func mapError(err error) httpserver.ErrorMapping {
	if mapping, ok := httpserver.MapRequestError(err); ok {
		return mapping
	}

	switch {
	case errors.Is(err, projecterrs.ErrProjectNotFound):
		return httpserver.NewErrorMapping(http.StatusNotFound, "project_not_found", "project not found", false)

	case errors.Is(err, projecterrs.ErrElementNotFound):
		return httpserver.NewErrorMapping(http.StatusNotFound, "element_not_found", "element not found", false)

	case errors.Is(err, projecterrs.ErrElementKeyAlreadyExists):
		return httpserver.NewErrorMapping(http.StatusConflict, "element_key_already_exists", "element key already exists", false)

	case errors.Is(err, projecterrs.ErrElementInUse):
		return httpserver.NewErrorMapping(http.StatusConflict, "element_in_use", "element is used by an active step", false)

	case errors.Is(err, projecterrs.ErrProjectIDRequired),
		errors.Is(err, projecterrs.ErrElementProjectIDRequired):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_project_id", "project ID is required", false)

	case errors.Is(err, projecterrs.ErrElementIDRequired):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_element_id", "element ID is required", false)

	case errors.Is(err, projecterrs.ErrLimitInvalid):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_limit", "limit must be between 1 and 100", false)

	case errors.Is(err, projecterrs.ErrOffsetInvalid):
		return httpserver.NewErrorMapping(http.StatusBadRequest, "invalid_offset", "offset must be greater than or equal to 0", false)

	case errors.Is(err, projecterrs.ErrProjectNameRequired):
		return httpserver.NewValidationErrorMapping("project_name_required", "project name is required", "name")

	case errors.Is(err, projecterrs.ErrProjectNameTooLong):
		return httpserver.NewValidationErrorMapping("project_name_too_long", "project name is too long", "name")

	case errors.Is(err, projecterrs.ErrElementKeyRequired):
		return httpserver.NewValidationErrorMapping("element_key_required", "element key is required", "key")

	case errors.Is(err, projecterrs.ErrElementKeyTooLong):
		return httpserver.NewValidationErrorMapping("element_key_too_long", "element key is too long", "key")

	case errors.Is(err, projecterrs.ErrElementLabelRequired):
		return httpserver.NewValidationErrorMapping("element_label_required", "element label is required", "label")

	case errors.Is(err, projecterrs.ErrElementLabelTooLong):
		return httpserver.NewValidationErrorMapping("element_label_too_long", "element label is too long", "label")

	case errors.Is(err, projecterrs.ErrElementDescriptionTooLong):
		return httpserver.NewValidationErrorMapping("element_description_too_long", "element description is too long", "description")

	case errors.Is(err, projecterrs.ErrEmptyElementUpdateParams):
		return httpserver.NewErrorMapping(http.StatusUnprocessableEntity, "empty_update", "at least one field must be provided", false)

	case errors.Is(err, projecterrs.ErrFailedGenerateUniqueKey):
		return httpserver.InternalErrorMapping()

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
			"http handler failed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Error(err),
		)
	}

	httpserver.WriteJSON(w, mapping.Status, mapping.Response)
}
