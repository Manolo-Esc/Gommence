package netw

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// use: err := Encode(w, r, http.StatusOK, obj)
func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}

// use: decoded, err := Decode[CreateSomethingRequest](r)
func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("error decoding json: %w", err)
	}
	return v, nil
}
