package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
)

const (
	projectKeyHeader = "X-Project-Key"
)

func ExtractProjectKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectKey := strings.TrimSpace(r.Header.Get(projectKeyHeader))
		if projectKey == "" {
			writeProjectKeyError(w, http.StatusForbidden, "project_key_required", "X-Project-Key header is required")
			return
		}
		ctx := context.WithValue(r.Context(), "projectKey", projectKey)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeProjectKeyError(w http.ResponseWriter, status int, code, message string) {
	httpserver.WriteJSON(w, status, httpserver.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
