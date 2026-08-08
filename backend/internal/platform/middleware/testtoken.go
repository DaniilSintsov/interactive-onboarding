package middleware

import (
	"context"
	"net/http"
	"strings"
)

const (
	testTokenHeader = "X-Scenario-Test-Token"
)

func ExtractTestToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testToken := strings.TrimSpace(r.Header.Get(testTokenHeader))
		ctx := context.WithValue(r.Context(), "testToken", testToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
