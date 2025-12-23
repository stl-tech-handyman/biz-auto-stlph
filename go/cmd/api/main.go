package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bizops360/go-api/internal/app"
	"github.com/bizops360/go-api/internal/config"
	"github.com/bizops360/go-api/internal/domain"
	httphandler "github.com/bizops360/go-api/internal/http"
	"github.com/bizops360/go-api/internal/infra/db"
	logger "github.com/bizops360/go-api/internal/infra/log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger := logger.NewLogger(cfg.LogLevel)

	logger.Info("starting server",
		"port", cfg.Port,
		"environment", cfg.Environment,
		"configDir", cfg.ConfigDir,
		"templatesDir", cfg.TemplatesDir,
	)

	// Initialize business loader
	businessLoader := config.NewBusinessLoader(cfg)

	// Initialize repositories (stub implementations for now)
	jobsRepo := db.NewMemoryJobsRepo()

	// Register pipeline actions (stub implementations)
	actions := map[string]domain.Action{
		"normalize_input":         &app.NormalizeInputAction{},
		"send_slack_notification": &app.SendSlackNotificationAction{},
		// Add more actions as they're implemented
	}

	// Initialize pipeline runner
	pipelineRunner := domain.NewPipelineRunner(actions)

	// Initialize services
	formEventsService := app.NewFormEventsService(businessLoader, pipelineRunner, jobsRepo)
	triggersService := app.NewTriggersService(businessLoader, pipelineRunner, jobsRepo)

	// Initialize router
	router := httphandler.NewRouter(formEventsService, triggersService, businessLoader, logger, cfg.Environment)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
