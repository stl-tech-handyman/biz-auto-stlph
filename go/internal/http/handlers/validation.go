package handlers

import (
	"net/http"
	"strconv"

	"github.com/bizops360/go-api/internal/util"
)

// ValidateMethod validates that the request method matches the expected method
func ValidateMethod(r *http.Request, expectedMethod string, w http.ResponseWriter) bool {
	if r.Method != expectedMethod {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return false
	}
	return true
}

// ValidateRequiredString validates that a string field is not empty
func ValidateRequiredString(value, fieldName string, w http.ResponseWriter) bool {
	if value == "" {
		util.WriteError(w, http.StatusBadRequest, fieldName+" is required")
		return false
	}
	return true
}

// ParseFloatFromQuery parses a float from query parameter
func ParseFloatFromQuery(r *http.Request, key string) (*float64, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return nil, nil
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// ParseFloatFromMap parses a float from a map[string]interface{}
func ParseFloatFromMap(m map[string]interface{}, key string) (*float64, error) {
	val, ok := m[key]
	if !ok {
		return nil, nil
	}
	switch v := val.(type) {
	case float64:
		return &v, nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		return &f, nil
	default:
		return nil, nil
	}
}

// GetStringFromMap safely gets a string value from a map
func GetStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// GetBoolFromMap safely gets a bool value from a map
func GetBoolFromMap(m map[string]interface{}, key string) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return false
}

