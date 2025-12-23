package handlers

import (
	"embed"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bizops360/go-api/internal/util"
	"gopkg.in/yaml.v3"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


//go:embed swagger-ui.html
var swaggerUIHTML embed.FS

// SwaggerHandler handles Swagger UI endpoint
type SwaggerHandler struct {
	openAPIPath string
}

// NewSwaggerHandler creates a new Swagger handler
func NewSwaggerHandler(openAPIPath string) *SwaggerHandler {
	return &SwaggerHandler{
		openAPIPath: openAPIPath,
	}
}

// HandleSwaggerUI serves the Swagger UI page
func (h *SwaggerHandler) HandleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	// #region agent log
	if logFile, err := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "B", "location": "swagger_handler.go:32", "message": "HandleSwaggerUI called", "data": map[string]interface{}{"path": r.URL.Path, "method": r.Method, "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion
	
	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Read Swagger UI HTML template
	tmplContent, err := swaggerUIHTML.ReadFile("swagger-ui.html")
	
	// #region agent log
	if logFile, err2 := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		}
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "C", "location": "swagger_handler.go:48", "message": "ReadFile swagger-ui.html result", "data": map[string]interface{}{"error": errorMsg, "contentLength": len(tmplContent), "hasContent": len(tmplContent) > 0, "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion
	
	if err != nil {
		// #region agent log
		if logFile, err2 := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "C", "location": "swagger_handler.go:58", "message": "ReadFile failed, returning error", "data": map[string]interface{}{"error": err.Error(), "timestamp": time.Now().UnixMilli()}})
			logFile.Close()
		}
		// #endregion
		util.WriteError(w, http.StatusInternalServerError, "failed to load Swagger UI template: "+err.Error())
		return
	}
	
	// #region agent log
	if logFile, err2 := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "C", "location": "swagger_handler.go:66", "message": "Writing Swagger UI HTML to response", "data": map[string]interface{}{"contentLength": len(tmplContent), "firstChars": string(tmplContent[:min(100, len(tmplContent))]), "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(tmplContent)
}

// HandleOpenAPISpec serves the OpenAPI spec as JSON
func (h *SwaggerHandler) HandleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	// #region agent log
	if logFile, err := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "B", "location": "swagger_handler.go:52", "message": "HandleOpenAPISpec called", "data": map[string]interface{}{"path": r.URL.Path, "method": r.Method, "openAPIPath": h.openAPIPath, "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion
	
	// Support GET, HEAD, and OPTIONS methods
	// #region agent log
	if logFile, err2 := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "F", "location": "swagger_handler.go:99", "message": "Method check", "data": map[string]interface{}{"method": r.Method, "methodGet": http.MethodGet, "methodHead": http.MethodHead, "methodOptions": http.MethodOptions, "isGet": r.Method == http.MethodGet, "isHead": r.Method == http.MethodHead, "isOptions": r.Method == http.MethodOptions, "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion
	
	if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
		// #region agent log
		if logFile, err2 := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "F", "location": "swagger_handler.go:110", "message": "Method not allowed", "data": map[string]interface{}{"method": r.Method, "timestamp": time.Now().UnixMilli()}})
			logFile.Close()
		}
		// #endregion
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	
	// Handle OPTIONS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	openAPIContent, err := os.ReadFile(h.openAPIPath)
	
	// #region agent log
	if logFile, err2 := os.OpenFile("../../.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "D", "location": "swagger_handler.go:60", "message": "ReadFile OpenAPI spec result", "data": map[string]interface{}{"error": func() string { if err != nil { return err.Error() }; return "" }(), "contentLength": len(openAPIContent), "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion
	
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to read OpenAPI spec: "+err.Error())
		return
	}

	// Convert YAML to JSON
	var specData interface{}
	if err := yaml.Unmarshal(openAPIContent, &specData); err != nil {
		// #region agent log
		if logFile, err2 := os.OpenFile(".cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "E", "location": "swagger_handler.go:70", "message": "YAML Unmarshal error", "data": map[string]interface{}{"error": err.Error(), "timestamp": time.Now().UnixMilli()}})
			logFile.Close()
		}
		// #endregion
		util.WriteError(w, http.StatusInternalServerError, "failed to parse OpenAPI YAML: "+err.Error())
		return
	}

	jsonData, err := json.Marshal(specData)
	if err != nil {
		// #region agent log
		if logFile, err2 := os.OpenFile(".cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "E", "location": "swagger_handler.go:78", "message": "JSON Marshal error", "data": map[string]interface{}{"error": err.Error(), "timestamp": time.Now().UnixMilli()}})
			logFile.Close()
		}
		// #endregion
		util.WriteError(w, http.StatusInternalServerError, "failed to convert to JSON: "+err.Error())
		return
	}

	// #region agent log
	if logFile, err2 := os.OpenFile("../../.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "E", "location": "swagger_handler.go:85", "message": "OpenAPI spec converted successfully", "data": map[string]interface{}{"jsonLength": len(jsonData), "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Api-Key")
	
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}


// GetOpenAPIPath tries to find the OpenAPI spec file
func GetOpenAPIPath() string {
	// Try common locations - check from current working directory
	possiblePaths := []string{
		"./docs/api/openapi-ru.yaml",
		"./go/docs/api/openapi-ru.yaml",
		"docs/api/openapi-ru.yaml",
		"go/docs/api/openapi-ru.yaml",
		"../docs/api/openapi-ru.yaml",
		"../../docs/api/openapi-ru.yaml",
	}

	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			absPath, err := filepath.Abs(path)
			// #region agent log
			if logFile, err2 := os.OpenFile(".cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
				json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "D", "location": "swagger_handler.go:96", "message": "OpenAPI file found", "data": map[string]interface{}{"path": path, "absPath": func() string { if err == nil { return absPath }; return path }(), "timestamp": time.Now().UnixMilli()}})
				logFile.Close()
			}
			// #endregion
			if err == nil {
				return absPath
			}
			return path
		}
		// #region agent log
		if logFile, err2 := os.OpenFile(".cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "D", "location": "swagger_handler.go:103", "message": "OpenAPI file not found at path", "data": map[string]interface{}{"path": path, "error": func() string { if _, err := os.Stat(path); err != nil { return err.Error() }; return "" }(), "timestamp": time.Now().UnixMilli()}})
			logFile.Close()
		}
		// #endregion
	}

	// Default fallback - relative to where server is run from
	// #region agent log
	if logFile, err := os.OpenFile(GetLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(logFile).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "D", "location": "swagger_handler.go:110", "message": "Using default OpenAPI path", "data": map[string]interface{}{"defaultPath": "./docs/api/openapi-ru.yaml", "timestamp": time.Now().UnixMilli()}})
		logFile.Close()
	}
	// #endregion
	return "./docs/api/openapi-ru.yaml"
}

