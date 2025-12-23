package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
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
// It first tries to load from .env file in the current directory, then falls back to system env vars
func LoadConfig() *Config {
	// #region agent log
	if logFile, err := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		wd, _ := os.Getwd()
		envFile := filepath.Join(wd, ".env")
		envFileExists := false
		if _, err := os.Stat(envFile); err == nil {
			envFileExists = true
		}
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "H1",
			"location":     "config.go:LoadConfig",
			"message":      "Before godotenv.Load - checking .env file",
			"data": map[string]interface{}{
				"workingDir":    wd,
				"envFilePath":   envFile,
				"envFileExists": envFileExists,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	// Try to load .env file (ignore error if file doesn't exist)
	err := godotenv.Load()
	// #region agent log
	if logFile, logErr := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); logErr == nil {
		gmailCreds := os.Getenv("GMAIL_CREDENTIALS_JSON")
		gmailFrom := os.Getenv("GMAIL_FROM")
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": "H1,H5",
			"location":     "config.go:LoadConfig",
			"message":      "After godotenv.Load - GMAIL env vars",
			"data": map[string]interface{}{
				"godotenvError":        err != nil,
				"godotenvErrorMsg":     func() string { if err != nil { return err.Error() } else { return "" } }(),
				"gmailCredsSet":        gmailCreds != "",
				"gmailCredsLength":     len(gmailCreds),
				"gmailCredsValue":      gmailCreds,
				"gmailFromSet":         gmailFrom != "",
				"gmailFromValue":       gmailFrom,
			},
			"timestamp": time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
		logFile.Close()
	}
	// #endregion
	
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
	// #region agent log
	if key == "PORT" {
		if logFile, err := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			rawValue := os.Getenv(key)
			containsBracket := false
			if rawValue != "" {
				containsBracket = (rawValue[0] == '[') || (len(rawValue) > 7 && rawValue[:7] == "[string")
			}
			logEntry := map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "A",
				"location":     "config.go:getEnv",
				"message":      "PORT env var raw value",
				"data": map[string]interface{}{
					"rawValue":       rawValue,
					"key":            key,
					"hasValue":       rawValue != "",
					"containsBracket": containsBracket,
					"length":         len(rawValue),
				},
				"timestamp": time.Now().UnixMilli(),
			}
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
	}
	// #endregion
	if value := os.Getenv(key); value != "" {
		// #region agent log
		if key == "PORT" {
			if logFile, err := os.OpenFile("../../.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				containsBracket := (value[0] == '[') || (len(value) > 7 && value[:7] == "[string")
				logEntry := map[string]interface{}{
					"sessionId":    "debug-session",
					"runId":        "run1",
					"hypothesisId": "B",
					"location":     "config.go:getEnv",
					"message":      "PORT env var returned",
					"data": map[string]interface{}{
						"value":          value,
						"length":         len(value),
						"containsBracket": containsBracket,
					},
					"timestamp": time.Now().UnixMilli(),
				}
				json.NewEncoder(logFile).Encode(logEntry)
				logFile.Close()
			}
		}
		// #endregion
		return value
	}
	// #region agent log
	if key == "PORT" {
		if logFile, err := os.OpenFile("c:\\Users\\Alexey\\Code\\biz-operating-system\\stlph\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			logEntry := map[string]interface{}{
				"sessionId":    "debug-session",
				"runId":        "run1",
				"hypothesisId": "D",
				"location":     "config.go:getEnv",
				"message":      "PORT using default",
				"data": map[string]interface{}{
					"defaultValue": defaultValue,
				},
				"timestamp": time.Now().UnixMilli(),
			}
			json.NewEncoder(logFile).Encode(logEntry)
			logFile.Close()
		}
	}
	// #endregion
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

