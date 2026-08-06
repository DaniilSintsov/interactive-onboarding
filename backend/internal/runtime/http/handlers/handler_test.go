package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	runtimeModel "github.com/DaniilSintsov/interactive-onboarding/internal/runtime/model"
)

type runtimeServiceMock struct {
	pageID string
	userID string
	result *runtimeModel.RuntimeScenario
	err    error
}

func (m *runtimeServiceMock) FindScenario(pageID string, userID string) (*runtimeModel.RuntimeScenario, error) {
	m.pageID = pageID
	m.userID = userID
	return m.result, m.err
}

func TestGetScenarioRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		statusCode int
	}{
		{name: "missing page", body: `{"user_id":"user-1"}`, statusCode: http.StatusUnprocessableEntity},
		{name: "empty page", body: `{"page":"","user_id":"user-1"}`, statusCode: http.StatusUnprocessableEntity},
		{name: "page too long", body: `{"page":"` + strings.Repeat("p", 2049) + `","user_id":"user-1"}`, statusCode: http.StatusUnprocessableEntity},
		{name: "missing user id", body: `{"page":"/home"}`, statusCode: http.StatusUnprocessableEntity},
		{name: "empty user id", body: `{"page":"/home","user_id":""}`, statusCode: http.StatusUnprocessableEntity},
		{name: "user id too long", body: `{"page":"/home","user_id":"` + strings.Repeat("u", 256) + `"}`, statusCode: http.StatusUnprocessableEntity},
		{name: "unknown field", body: `{"page":"/home","user_id":"user-1","unexpected":true}`, statusCode: http.StatusBadRequest},
		{name: "multiple JSON values", body: `{"page":"/home","user_id":"user-1"} {}`, statusCode: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/sdk/scenarios/resolve", strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			NewHandler(nil).GetScenario(response, req)

			if response.Code != tt.statusCode {
				t.Fatalf("expected status %d, got %d", tt.statusCode, response.Code)
			}
		})
	}
}

func TestGetScenarioAcceptsValidRequest(t *testing.T) {
	service := &runtimeServiceMock{
		result: &runtimeModel.RuntimeScenario{ID: "scenario-1", Name: "Welcome"},
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sdk/scenarios/resolve", strings.NewReader(`{"page":"/home","user_id":"user-1"}`))
	response := httptest.NewRecorder()

	NewHandler(service).GetScenario(response, req)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if service.pageID != "/home" || service.userID != "user-1" {
		t.Fatalf("unexpected service arguments: page=%q user_id=%q", service.pageID, service.userID)
	}

	var scenario runtimeModel.RuntimeScenario
	if err := json.NewDecoder(response.Body).Decode(&scenario); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if scenario.ID != "scenario-1" || scenario.Name != "Welcome" {
		t.Fatalf("unexpected scenario response: %+v", scenario)
	}
}

func TestGetScenarioReturnsNoContentWhenNoScenarioMatches(t *testing.T) {
	service := &runtimeServiceMock{}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sdk/scenarios/resolve", strings.NewReader(`{"page":"/home","user_id":"user-1"}`))
	response := httptest.NewRecorder()

	NewHandler(service).GetScenario(response, req)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
}

func TestRegisterRoutesRegistersResolveEndpoint(t *testing.T) {
	service := &runtimeServiceMock{result: &runtimeModel.RuntimeScenario{ID: "scenario-1"}}
	router := http.NewServeMux()
	NewHandler(service).RegisterRoutes(router)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v1/sdk/scenarios/resolve", strings.NewReader(`{"page":"/home","user_id":"user-1"}`)))

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
}

func TestGetScenarioPreservesFrontendDataAsJSONObject(t *testing.T) {
	service := &runtimeServiceMock{result: &runtimeModel.RuntimeScenario{
		ID: "scenario-1",
		Steps: []runtimeModel.RuntimeStep{{
			ID:           "step-1",
			FrontendData: json.RawMessage(`{"placement":"bottom"}`),
		}},
	}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sdk/scenarios/resolve", strings.NewReader(`{"page":"/home","user_id":"user-1"}`))
	response := httptest.NewRecorder()

	NewHandler(service).GetScenario(response, req)

	var payload struct {
		Steps []struct {
			FrontendData json.RawMessage `json:"frontend_data"`
		} `json:"steps"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got := string(payload.Steps[0].FrontendData); got != `{"placement":"bottom"}` {
		t.Fatalf("frontend_data = %s, want JSON object", got)
	}
}
