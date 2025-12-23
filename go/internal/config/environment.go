package config

import (
	"os"
)

// Environment represents the deployment environment
type Environment string

const (
	EnvDev  Environment = "dev"
	EnvProd Environment = "prod"
)

// GetEnvironment returns the current environment
func GetEnvironment() Environment {
	env := os.Getenv("ENV")
	if env == "prod" || env == "production" {
		return EnvProd
	}
	return EnvDev // Default to dev
}

// IsProduction returns true if running in production
func IsProduction() bool {
	return GetEnvironment() == EnvProd
}

// IsDevelopment returns true if running in development
func IsDevelopment() bool {
	return GetEnvironment() == EnvDev
}

// GetProjectID returns the GCP project ID for the current environment
func GetProjectID() string {
	if IsProduction() {
		return os.Getenv("GCP_PROJECT_ID_PROD")
	}
	return os.Getenv("GCP_PROJECT_ID_DEV")
}

// GetServiceName returns the Cloud Run service name for the current environment
func GetServiceName() string {
	if IsProduction() {
		return "bizops360-api-go-prod"
	}
	return "bizops360-api-go-dev"
}

