package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, payload any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("The body is empty")
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(payload); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("request body must contain a single JSON value")
		}
		return err
	}
	return nil
}

func raiseError(w http.ResponseWriter, errMsg string, err error, status int) {
	if writeErr := WriteJson(w, status, ErrorResponse{Code: http.StatusText(status), Message: errMsg}); writeErr != nil {
		log.Printf("failed to write error response: %v", writeErr)
	}
	log.Printf(errMsg+" %v", err)
}

// ErrorResponse mirrors the common OpenAPI error schema.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
