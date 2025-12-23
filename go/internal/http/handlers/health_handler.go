package handlers

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/bizops360/go-api/internal/util"
)

var startTime = time.Now()

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HandleHealth handles GET /api/health
func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	svcKey := os.Getenv("SERVICE_API_KEY")
	
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"service":     "bizops360-api-go",
		"environment": env,
		"version":     "1.0.0",
		"timestamp":   time.Now().Format(time.RFC3339),
		"uptime":      time.Since(startTime).Seconds(),
		"memory": map[string]interface{}{
			"usedMb":  m.Alloc / 1024 / 1024,
			"totalMb": m.Sys / 1024 / 1024,
		},
		"checks": map[string]interface{}{
			"runtime":     "healthy",
			"environment": "healthy",
			"timestamp":   "healthy",
		},
		"debug": map[string]interface{}{
			"serviceApiKeySet": svcKey != "",
			"serviceApiKeyLength": len(svcKey),
		},
	})
}

// HandleReady handles GET /api/health/ready
func (h *HealthHandler) HandleReady(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now().Format(time.RFC3339),
		"checks": map[string]interface{}{
			"server": "ok",
			"memory": m.Alloc < 200*1024*1024,
			"uptime": time.Since(startTime).Seconds() > 0,
		},
	})
}

// HandleLive handles GET /api/health/live
func (h *HealthHandler) HandleLive(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(startTime).Seconds(),
	})
}

func getEnv(key, defaultValue string) string {
	// This is a simple fallback - actual env vars are read via os.Getenv in config package
	return defaultValue
}
