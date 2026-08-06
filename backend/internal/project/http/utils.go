package projecthttp

import (
	"encoding/json"
	"errors"
	"fmt"
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
			return 0, 0, fmt.Errorf("invalid limit: %w", err)
		}

		limit = parsedLimit
	}

	if rawOffset := query.Get("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid offset: %w", err)
		}

		offset = parsedOffset
	}

	if limit < 1 || limit > 100 {
		return 0, 0, fmt.Errorf("invalid limit: limit=%d: limit must be between 1 and 100", limit)
	}

	if offset < 0 {
		return 0, 0, fmt.Errorf("invalid offset: offset=%d: offset must be >= 0", offset)
	}

	return limit, offset, nil
}

func parseUUIDPath(r *http.Request, name string) (uuid.UUID, error) {
	value := r.PathValue(name)

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse path uuid parameter %s: %w", name, err)
	}

	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%s must not be nil UUID", name)
	}

	return id, nil
}

func parseJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON value")
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
