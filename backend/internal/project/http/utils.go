package projecthttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

func parsePagination(r *http.Request) (int, int, error) {
	limit := defaultLimit
	offset := defaultOffset

	query := r.URL.Query()

	if rawLimit := query.Get("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil {
			return 0, 0, newRequestError(
				"invalid_limit",
				"limit must be an integer",
				map[string]any{"parameter": "limit"},
				err,
			)
		}

		limit = parsedLimit
	}

	if rawOffset := query.Get("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil {
			return 0, 0, newRequestError(
				"invalid_offset",
				"offset must be an integer",
				map[string]any{"parameter": "offset"},
				err,
			)
		}

		offset = parsedOffset
	}

	if limit < 1 || limit > 100 {
		return 0, 0, newRequestError(
			"invalid_limit",
			"limit must be between 1 and 100",
			map[string]any{"parameter": "limit"},
			nil,
		)
	}

	if offset < 0 {
		return 0, 0, newRequestError(
			"invalid_offset",
			"offset must be greater than or equal to 0",
			map[string]any{"parameter": "offset"},
			nil,
		)
	}

	return limit, offset, nil
}

func parseUUIDPath(r *http.Request, name string) (uuid.UUID, error) {
	value := r.PathValue(name)

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, newRequestError(
			pathParameterErrorCode(name),
			name+" must be a valid UUID",
			map[string]any{"parameter": name},
			err,
		)
	}

	if id == uuid.Nil {
		return uuid.Nil, newRequestError(
			pathParameterErrorCode(name),
			name+" must not be a nil UUID",
			map[string]any{"parameter": name},
			nil,
		)
	}

	return id, nil
}

func pathParameterErrorCode(name string) string {
	switch name {
	case "projectId":
		return "invalid_project_id"
	case "elementId":
		return "invalid_element_id"
	default:
		return "invalid_path_parameter"
	}
}

func parseJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return newRequestError(
			"invalid_request_body",
			"request body must contain valid JSON",
			map[string]any{"field": "body"},
			err,
		)
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return newRequestError(
			"invalid_request_body",
			"request body must contain exactly one JSON value",
			map[string]any{"field": "body"},
			err,
		)
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if value != nil {
		_ = json.NewEncoder(w).Encode(value)
	}
}
