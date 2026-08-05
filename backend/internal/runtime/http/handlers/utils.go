package handlers

import (
	"encoding/json"
	"fmt"
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
	return json.NewDecoder(r.Body).Decode(payload)
}

func raiseError(w http.ResponseWriter, errMsg string, err error, status int) {
	http.Error(w, "Something went wrong", status)
	log.Printf(errMsg+" %v", err)
}
