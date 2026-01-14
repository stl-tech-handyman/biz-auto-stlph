package handlers

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
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
	
	healthData := map[string]interface{}{
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
			"serviceApiKeySet":    svcKey != "",
			"serviceApiKeyLength": len(svcKey),
		},
	}

	// Check if client wants JSON (for API clients)
	acceptHeader := r.Header.Get("Accept")
	wantsJSON := strings.Contains(acceptHeader, "application/json") || r.URL.Query().Get("format") == "json"

	if wantsJSON {
		util.WriteJSON(w, http.StatusOK, healthData)
		return
	}

	// Otherwise, return friendly HTML page
	h.renderHealthHTML(w, healthData)
}

// renderHealthHTML renders a friendly HTML page with health information
func (h *HealthHandler) renderHealthHTML(w http.ResponseWriter, data map[string]interface{}) {
	service := data["service"].(string)
	env := data["environment"].(string)
	version := data["version"].(string)
	timestamp := data["timestamp"].(string)
	uptime := data["uptime"].(float64)
	
	memory := data["memory"].(map[string]interface{})
	usedMb := memory["usedMb"].(uint64)
	totalMb := memory["totalMb"].(uint64)
	var memoryPercent float64
	if totalMb > 0 {
		memoryPercent = float64(usedMb) / float64(totalMb) * 100
	}
	
	checks := data["checks"].(map[string]interface{})
	debug := data["debug"].(map[string]interface{})
	
	// Format uptime
	uptimeHours := int(uptime) / 3600
	uptimeMinutes := (int(uptime) % 3600) / 60
	uptimeSeconds := int(uptime) % 60
	uptimeStr := fmt.Sprintf("%dh %dm %ds", uptimeHours, uptimeMinutes, uptimeSeconds)
	
	// Determine environment badge color
	envBadgeColor := "bg-info"
	if env == "prod" {
		envBadgeColor = "bg-success"
	} else if env == "dev" {
		envBadgeColor = "bg-warning"
	}
	
	// Determine API key status
	apiKeySet := debug["serviceApiKeySet"].(bool)
	apiKeyIcon := "check-circle-fill"
	apiKeyColor := "#28a745"
	apiKeyStatus := "Set"
	if !apiKeySet {
		apiKeyIcon = "x-circle-fill"
		apiKeyColor = "#dc3545"
		apiKeyStatus = "Not Set"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>System Health - %s</title>
    <link href="/static/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/css/bootstrap-icons.css" rel="stylesheet">
    <style>
        body {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
            padding: 2rem 0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
        }
        .health-container {
            max-width: 900px;
            margin: 0 auto;
        }
        .health-card {
            background: white;
            border-radius: 16px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .health-header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 2rem;
            text-align: center;
        }
        .health-header h1 {
            margin: 0;
            font-size: 2rem;
            font-weight: 600;
        }
        .health-header .badge {
            margin-top: 0.5rem;
            font-size: 0.9rem;
            padding: 0.5rem 1rem;
        }
        .health-body {
            padding: 2rem;
        }
        .stat-card {
            background: #f8f9fa;
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 1rem;
            border-left: 4px solid #667eea;
        }
        .stat-label {
            font-size: 0.875rem;
            color: #6c757d;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 0.5rem;
        }
        .stat-value {
            font-size: 1.75rem;
            font-weight: 600;
            color: #212529;
        }
        .check-item {
            display: flex;
            align-items: center;
            padding: 0.75rem;
            background: #f8f9fa;
            border-radius: 8px;
            margin-bottom: 0.5rem;
        }
        .check-item i {
            margin-right: 0.75rem;
            font-size: 1.25rem;
        }
        .check-item.healthy i {
            color: #28a745;
        }
        .memory-bar {
            height: 8px;
            background: #e9ecef;
            border-radius: 4px;
            overflow: hidden;
            margin-top: 0.5rem;
        }
        .memory-bar-fill {
            height: 100%%;
            background: linear-gradient(90deg, #28a745 0%%, #20c997 100%%);
            transition: width 0.3s ease;
        }
        .json-link {
            text-align: center;
            margin-top: 2rem;
            padding-top: 2rem;
            border-top: 1px solid #dee2e6;
        }
        .timestamp {
            color: #6c757d;
            font-size: 0.875rem;
        }
    </style>
</head>
<body>
    <div class="health-container">
        <div class="health-card">
            <div class="health-header">
                <h1><i class="bi bi-heart-pulse-fill me-2"></i>System Health</h1>
                <span class="badge %s">%s</span>
            </div>
            <div class="health-body">
                <!-- Service Info -->
                <div class="stat-card">
                    <div class="stat-label">Service</div>
                    <div class="stat-value">%s</div>
                    <div class="timestamp mt-2">
                        <i class="bi bi-tag-fill me-1"></i>Version %s
                    </div>
                </div>

                <!-- Status Checks -->
                <div class="stat-card">
                    <div class="stat-label mb-3">Status Checks</div>
                    <div class="check-item healthy">
                        <i class="bi bi-check-circle-fill"></i>
                        <div>
                            <strong>Runtime:</strong> %s
                        </div>
                    </div>
                    <div class="check-item healthy">
                        <i class="bi bi-check-circle-fill"></i>
                        <div>
                            <strong>Environment:</strong> %s
                        </div>
                    </div>
                    <div class="check-item healthy">
                        <i class="bi bi-check-circle-fill"></i>
                        <div>
                            <strong>Timestamp:</strong> %s
                        </div>
                    </div>
                </div>

                <!-- Memory Usage -->
                <div class="stat-card">
                    <div class="stat-label">Memory Usage</div>
                    <div class="stat-value">%d MB / %d MB</div>
                    <div class="memory-bar">
                        <div class="memory-bar-fill" style="width: %.1f%%"></div>
                    </div>
                    <div class="timestamp mt-2">
                        %.1f%% of total memory used
                    </div>
                </div>

                <!-- Uptime -->
                <div class="stat-card">
                    <div class="stat-label">Uptime</div>
                    <div class="stat-value">%s</div>
                    <div class="timestamp mt-2">
                        <i class="bi bi-clock-fill me-1"></i>Server has been running since startup
                    </div>
                </div>

                <!-- Debug Info -->
                <div class="stat-card">
                    <div class="stat-label">Configuration</div>
                    <div class="check-item">
                        <i class="bi bi-%s" style="color: %s;"></i>
                        <div>
                            <strong>Service API Key:</strong> %s (Length: %d)
                        </div>
                    </div>
                </div>

                <!-- Timestamp -->
                <div class="text-center timestamp mt-3">
                    <i class="bi bi-calendar-event me-1"></i>Last checked: %s
                </div>

                <!-- JSON Link -->
                <div class="json-link">
                    <a href="?format=json" class="btn btn-outline-primary">
                        <i class="bi bi-code-slash me-2"></i>View as JSON
                    </a>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`,
		service,
		envBadgeColor, env,
		service, version,
		checks["runtime"], checks["environment"], checks["timestamp"],
		usedMb, totalMb, memoryPercent, memoryPercent,
		uptimeStr,
		apiKeyIcon, apiKeyColor,
		apiKeyStatus,
		debug["serviceApiKeyLength"].(int),
		timestamp,
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
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

// HandleTime handles GET /api/time - returns current server time
func (h *HealthHandler) HandleTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	
	now := time.Now()
	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"timestamp": now.Format(time.RFC3339),
		"unix":      now.Unix(),
		"unixMilli": now.UnixMilli(),
		"time":      now.Format("15:04:05"),
		"date":      now.Format("2006-01-02"),
	})
}

func getEnv(key, defaultValue string) string {
	// This is a simple fallback - actual env vars are read via os.Getenv in config package
	return defaultValue
}
