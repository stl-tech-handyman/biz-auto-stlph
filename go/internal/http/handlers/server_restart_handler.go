package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/bizops360/go-api/internal/util"
)

// ServerRestartHandler handles server restart requests
type ServerRestartHandler struct {
	logger *slog.Logger
}

// NewServerRestartHandler creates a new server restart handler
func NewServerRestartHandler(logger *slog.Logger) *ServerRestartHandler {
	return &ServerRestartHandler{
		logger: logger,
	}
}

// HandleRestart handles POST /api/server/restart
func (h *ServerRestartHandler) HandleRestart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Only allow in development environment
	env := os.Getenv("ENV")
	if env != "dev" && env != "" {
		util.WriteError(w, http.StatusForbidden, "server restart only allowed in development environment")
		return
	}

	h.logger.Info("Server restart requested via API")

	// Get the directory of the executable
	execPath, err := os.Executable()
	if err != nil {
		h.logger.Error("Failed to get executable path", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "failed to determine server path")
		return
	}

	execDir := filepath.Dir(execPath)
	
	// Try to find restart script in common locations
	var restartScript string
	possibleScripts := []string{
		filepath.Join(execDir, "restart-server.sh"),
		filepath.Join(execDir, "restart-server.bat"),
		filepath.Join(execDir, "..", "restart-server.sh"),
		filepath.Join(execDir, "..", "restart-server.bat"),
		"./restart-server.sh",
		"./restart-server.bat",
	}

	for _, script := range possibleScripts {
		if _, err := os.Stat(script); err == nil {
			restartScript = script
			break
		}
	}

	if restartScript == "" {
		h.logger.Warn("Restart script not found, attempting direct restart")
		// Fallback: try to restart by killing the process
		// This is a simple approach - in production you'd want a proper process manager
		go func() {
			time.Sleep(1 * time.Second) // Give time for response to be sent
			os.Exit(0) // Exit and let process manager restart
		}()
		
		util.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"status":  "restarting",
			"message": "Server restart initiated",
		})
		return
	}

	// Execute restart script in background
	go func() {
		time.Sleep(500 * time.Millisecond) // Give time for HTTP response

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			// On Windows, use cmd.exe to run .bat file
			if filepath.Ext(restartScript) == ".bat" {
				cmd = exec.Command("cmd.exe", "/c", restartScript)
			} else {
				// Try with bash (Git Bash)
				cmd = exec.Command("bash", restartScript)
			}
		} else {
			cmd = exec.Command("bash", restartScript)
		}

		cmd.Dir = filepath.Dir(restartScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		h.logger.Info("Executing restart script", "script", restartScript)
		if err := cmd.Start(); err != nil {
			h.logger.Error("Failed to start restart script", "error", err)
		}
	}()

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "restarting",
		"message": "Server restart initiated",
		"script":  restartScript,
	})
}
