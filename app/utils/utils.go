package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"chatbot/types"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// Decode JSON from struct
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

// Encode into JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Cache-Control", "no-cache")

	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Encode error message into JSON
func WriteError(w http.ResponseWriter, status int, err error) {
	responseFormat := types.APIResponse{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}
	WriteJSON(w, status, responseFormat)
}

func CurlRequest(url string, payload map[string]interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Check for non-200 responses here if needed
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return resp, nil // Return the *http.Response directly for streaming
}
