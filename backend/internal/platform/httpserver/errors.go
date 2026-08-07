package httpserver

import (
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type ErrorMapping struct {
	Status   int
	Response ErrorResponse
	Log      bool
}

type RequestError struct {
	code    string
	message string
	details map[string]any
	cause   error
}

func (e *RequestError) Error() string {
	if e.cause == nil {
		return e.message
	}

	return e.message + ": " + e.cause.Error()
}

func (e *RequestError) Unwrap() error {
	return e.cause
}

func NewRequestError(
	code string,
	message string,
	details map[string]any,
	cause error,
) error {
	return &RequestError{
		code:    code,
		message: message,
		details: details,
		cause:   cause,
	}
}

func MapRequestError(err error) (ErrorMapping, bool) {
	var requestErr *RequestError
	if !errors.As(err, &requestErr) {
		return ErrorMapping{}, false
	}

	return ErrorMapping{
		Status: http.StatusBadRequest,
		Response: ErrorResponse{
			Code:    requestErr.code,
			Message: requestErr.message,
			Details: requestErr.details,
		},
	}, true
}

func NewErrorMapping(status int, code, message string, logError bool) ErrorMapping {
	return ErrorMapping{
		Status: status,
		Response: ErrorResponse{
			Code:    code,
			Message: message,
		},
		Log: logError,
	}
}

func NewValidationErrorMapping(code, message, field string) ErrorMapping {
	mapping := NewErrorMapping(http.StatusUnprocessableEntity, code, message, false)
	mapping.Response.Details = map[string]any{"field": field}

	return mapping
}

func InternalErrorMapping() ErrorMapping {
	return NewErrorMapping(
		http.StatusInternalServerError,
		"internal_error",
		"internal server error",
		true,
	)
}
