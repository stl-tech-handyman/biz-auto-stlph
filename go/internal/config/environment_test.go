package config

import (
	"os"
	"testing"
)

func TestGetEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected Environment
	}{
		{
			name:     "dev environment",
			envValue: "dev",
			expected: EnvDev,
		},
		{
			name:     "prod environment",
			envValue: "prod",
			expected: EnvProd,
		},
		{
			name:     "production environment",
			envValue: "production",
			expected: EnvProd,
		},
		{
			name:     "default to dev",
			envValue: "",
			expected: EnvDev,
		},
		{
			name:     "unknown defaults to dev",
			envValue: "unknown",
			expected: EnvDev,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			original := os.Getenv("ENV")
			defer os.Setenv("ENV", original)

			// Set test value
			if tt.envValue != "" {
				os.Setenv("ENV", tt.envValue)
			} else {
				os.Unsetenv("ENV")
			}

			result := GetEnvironment()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsProduction(t *testing.T) {
	original := os.Getenv("ENV")
	defer os.Setenv("ENV", original)

	os.Setenv("ENV", "prod")
	if !IsProduction() {
		t.Error("expected IsProduction() to return true when ENV=prod")
	}

	os.Setenv("ENV", "dev")
	if IsProduction() {
		t.Error("expected IsProduction() to return false when ENV=dev")
	}
}

func TestIsDevelopment(t *testing.T) {
	original := os.Getenv("ENV")
	defer os.Setenv("ENV", original)

	os.Setenv("ENV", "dev")
	if !IsDevelopment() {
		t.Error("expected IsDevelopment() to return true when ENV=dev")
	}

	os.Setenv("ENV", "prod")
	if IsDevelopment() {
		t.Error("expected IsDevelopment() to return false when ENV=prod")
	}
}

func TestGetServiceName(t *testing.T) {
	original := os.Getenv("ENV")
	defer os.Setenv("ENV", original)

	os.Setenv("ENV", "dev")
	if GetServiceName() != "bizops360-api-go-dev" {
		t.Errorf("expected 'bizops360-api-go-dev', got '%s'", GetServiceName())
	}

	os.Setenv("ENV", "prod")
	if GetServiceName() != "bizops360-api-go-prod" {
		t.Errorf("expected 'bizops360-api-go-prod', got '%s'", GetServiceName())
	}
}

