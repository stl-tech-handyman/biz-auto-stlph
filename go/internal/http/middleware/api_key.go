package middleware

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/util"
)

// APIKeyMiddleware validates X-Api-Key header
func APIKeyMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := util.GetRequestID(r.Context())
		
		svcKey := os.Getenv("SERVICE_API_KEY")
		
		// CRITICAL DEBUG: Write to stderr to confirm middleware is called
		fmt.Fprintf(os.Stderr, "[DEBUG] APIKeyMiddleware CALLED: path=%s, svcKeySet=%v, svcKeyLength=%d\n", r.URL.Path, svcKey != "", len(svcKey))
		
		// #region agent log
		logPath := func() string {
			if cwd, err := os.Getwd(); err == nil {
				// Server runs from go/ directory, so go up one level
				if filepath.Base(cwd) == "go" {
					return filepath.Join(cwd, "..", ".cursor", "debug.log")
				}
				return filepath.Join(cwd, ".cursor", "debug.log")
			}
			return "../../.cursor/debug.log"
		}()
		if logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			cwd := "unknown"
			if c, err := os.Getwd(); err == nil {
				cwd = c
			}
			json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "G", "location": "api_key.go:17", "message": "APIKeyMiddleware ENTRY - SERVICE_API_KEY check", "data": map[string]interface{}{"svcKeySet": svcKey != "", "svcKeyLength": len(svcKey), "svcKeyPrefix": func() string { if len(svcKey) > 0 { return svcKey[:min(10, len(svcKey))] + "..." }; return "" }(), "path": r.URL.Path, "method": r.Method, "cwd": cwd, "logPath": logPath, "timestamp": time.Now().UnixMilli()}})
			logFile.Close()
		}
		// #endregion
		
		// Log API key check details for debugging
		logger.Info("api_key_check",
			"requestId", requestID,
			"path", r.URL.Path,
			"method", r.Method,
			"svcKeySet", svcKey != "",
			"svcKeyLength", len(svcKey),
			"svcKeyPrefix", func() string {
				if len(svcKey) > 0 {
					return svcKey[:min(10, len(svcKey))] + "..."
				}
				return ""
			}(),
		)
		
		if svcKey == "" {
			// #region agent log
			logPath := func() string {
				if cwd, err := os.Getwd(); err == nil {
					// Server runs from go/ directory, so go up one level
					if filepath.Base(cwd) == "go" {
						return filepath.Join(cwd, "..", ".cursor", "debug.log")
					}
					return filepath.Join(cwd, ".cursor", "debug.log")
				}
				return "../../.cursor/debug.log"
			}()
			if logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				envVars := []string{}
				for _, v := range os.Environ() {
					if strings.Contains(v, "SERVICE") || strings.Contains(v, "API") {
						if len(v) > 50 {
							envVars = append(envVars, v[:50])
						} else {
							envVars = append(envVars, v)
						}
					}
				}
				json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "G", "location": "api_key.go:50", "message": "SERVICE_API_KEY is empty", "data": map[string]interface{}{"path": r.URL.Path, "allEnvVars": envVars, "cwd": func() string { if cwd, err := os.Getwd(); err == nil { return cwd }; return "unknown" }(), "logPath": logPath, "timestamp": time.Now().UnixMilli()}})
				logFile.Close()
			}
			// #endregion
			logger.Error("api_key_not_configured",
				"requestId", requestID,
				"path", r.URL.Path,
			)
			http.Error(w, `{"error":"Service Configuration Error","message":"API authentication is not properly configured"}`, http.StatusInternalServerError)
			return
		}

		apiKey := r.Header.Get("X-Api-Key")
		if apiKey == "" {
			apiKey = r.Header.Get("x-api-key")
		}

		if apiKey == "" {
			logger.Warn("api_key_missing",
				"requestId", requestID,
				"path", r.URL.Path,
				"ip", r.RemoteAddr,
			)
			http.Error(w, `{"error":"Unauthorized","message":"API key is required","hint":"Include X-Api-Key header in your request"}`, http.StatusUnauthorized)
			return
		}

		trimmedApiKey := strings.TrimSpace(apiKey)
		trimmedSvcKey := strings.TrimSpace(svcKey)
		
		// Log comparison details
		logger.Info("api_key_comparison",
			"requestId", requestID,
			"path", r.URL.Path,
			"apiKeyLength", len(apiKey),
			"apiKeyPrefix", func() string {
				if len(apiKey) > 0 {
					return apiKey[:min(10, len(apiKey))] + "..."
				}
				return ""
			}(),
			"trimmedApiKeyLength", len(trimmedApiKey),
			"trimmedSvcKeyLength", len(trimmedSvcKey),
			"keysMatch", trimmedApiKey == trimmedSvcKey,
		)

		if trimmedApiKey != trimmedSvcKey {
			logger.Warn("api_key_invalid",
				"requestId", requestID,
				"path", r.URL.Path,
				"ip", r.RemoteAddr,
				"receivedKeyLength", len(apiKey),
				"expectedKeyLength", len(svcKey),
				"receivedKeyPrefix", func() string {
					if len(apiKey) > 0 {
						return apiKey[:min(10, len(apiKey))] + "..."
					}
					return ""
				}(),
				"expectedKeyPrefix", func() string {
					if len(svcKey) > 0 {
						return svcKey[:min(10, len(svcKey))] + "..."
					}
					return ""
				}(),
			)
			http.Error(w, `{"error":"Unauthorized","message":"Invalid API key"}`, http.StatusUnauthorized)
			return
		}

		logger.Info("api_key_valid",
			"requestId", requestID,
			"path", r.URL.Path,
		)
		
		next.ServeHTTP(w, r)
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

