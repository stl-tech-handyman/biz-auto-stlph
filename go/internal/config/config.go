package config

import (
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	Port            string
	Environment    string
	ConfigDir       string
	TemplatesDir    string
	LogLevel        string
	IsProduction    bool
	IsDevelopment   bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	env := GetEnvironment()
	isProd := IsProduction()
	isDev := IsDevelopment()
	
	// Set default log level based on environment
	defaultLogLevel := "info"
	if isDev {
		defaultLogLevel = "debug"
	}
	
	return &Config{
		Port:          getEnv("PORT", "8080"),
		Environment:   string(env),
		ConfigDir:     getEnv("CONFIG_DIR", "/app/config"),
		TemplatesDir:  getEnv("TEMPLATES_DIR", "/app/templates"),
		LogLevel:      getEnv("LOG_LEVEL", defaultLogLevel),
		IsProduction:  isProd,
		IsDevelopment: isDev,
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetBusinessConfigPath returns the path to a business config file
func (c *Config) GetBusinessConfigPath(businessID string) string {
	return filepath.Join(c.ConfigDir, "businesses", businessID+".yaml")
}

// GetPipelineConfigPath returns the path to a pipeline config file
func (c *Config) GetPipelineConfigPath(pipelineKey string) string {
	return filepath.Join(c.ConfigDir, "pipelines", pipelineKey+".yaml")
}

// GetTemplatePath returns the path to a template file
func (c *Config) GetTemplatePath(businessID, templateName string) string {
	return filepath.Join(c.TemplatesDir, businessID, templateName)
}

