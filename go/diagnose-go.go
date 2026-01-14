package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	logPath := `c:\Users\Alexey\Code\biz-operating-system\stlph\.cursor\debug.log`
	
	log := func(hypothesisId, location, message string, data map[string]interface{}) {
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer logFile.Close()
		
		logEntry := map[string]interface{}{
			"sessionId":    "debug-session",
			"runId":        "run1",
			"hypothesisId": hypothesisId,
			"location":     location,
			"message":      message,
			"data":         data,
			"timestamp":    time.Now().UnixMilli(),
		}
		json.NewEncoder(logFile).Encode(logEntry)
	}

	// Hypothesis A: Check Go version
	goVersion, _ := exec.Command("go", "version").Output()
	log("A", "diagnose-go.go:go-version", "Go version check", map[string]interface{}{
		"version": string(goVersion),
		"raw":     fmt.Sprintf("%q", goVersion),
	})

	// Hypothesis B: Check GOROOT
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		goroot = runtime.GOROOT()
	}
	log("B", "diagnose-go.go:goroot", "GOROOT check", map[string]interface{}{
		"goroot":        goroot,
		"gorootExists": dirExists(goroot),
		"stdlibPath":    filepath.Join(goroot, "src"),
		"stdlibExists":  dirExists(filepath.Join(goroot, "src")),
	})

	// Hypothesis C: Check GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
		if runtime.GOOS == "windows" {
			gopath = filepath.Join(os.Getenv("USERPROFILE"), "go")
		}
	}
	log("C", "diagnose-go.go:gopath", "GOPATH check", map[string]interface{}{
		"gopath":       gopath,
		"gopathExists": dirExists(gopath),
	})

	// Hypothesis D: Check toolchain directory
	toolchainDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "go", "toolchain")
	if runtime.GOOS != "windows" {
		toolchainDir = filepath.Join(os.Getenv("HOME"), ".cache", "go-build")
	}
	log("D", "diagnose-go.go:toolchain", "Toolchain directory check", map[string]interface{}{
		"toolchainDir":       toolchainDir,
		"toolchainDirExists": dirExists(toolchainDir),
	})

	// Hypothesis E: Try to run go env
	goEnv, _ := exec.Command("go", "env", "GOROOT", "GOPATH", "GOTOOLCHAIN").Output()
	log("E", "diagnose-go.go:go-env", "Go environment variables", map[string]interface{}{
		"goEnvOutput": string(goEnv),
	})

	// Hypothesis F: Check if go.mod requires different version
	goModPath := "go.mod"
	if _, err := os.Stat(goModPath); err == nil {
		goModContent, _ := os.ReadFile(goModPath)
		log("F", "diagnose-go.go:go-mod", "go.mod content check", map[string]interface{}{
			"goModExists": true,
			"goModSize":  len(goModContent),
			"firstLines": string(goModContent[:min(200, len(goModContent))]),
		})
	} else {
		log("F", "diagnose-go.go:go-mod", "go.mod content check", map[string]interface{}{
			"goModExists": false,
			"error":       err.Error(),
		})
	}

	// Hypothesis G: Try go list to see what error we get
	goListCmd := exec.Command("go", "list", "-m")
	goListCmd.Dir = "."
	goListOutput, goListErr := goListCmd.CombinedOutput()
	log("G", "diagnose-go.go:go-list", "go list command test", map[string]interface{}{
		"output": string(goListOutput),
		"error":  fmt.Sprintf("%v", goListErr),
		"exitCode": getExitCode(goListErr),
	})

	fmt.Println("Diagnostics complete. Check debug.log for details.")
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return -1
}
