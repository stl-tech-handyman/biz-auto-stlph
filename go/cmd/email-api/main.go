package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bizops360/go-api/internal/config"
	emailapi "github.com/bizops360/go-api/internal/http/emailapi"
	"github.com/bizops360/go-api/internal/infra/email"
	logger "github.com/bizops360/go-api/internal/infra/log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	log := logger.NewLogger(cfg.LogLevel)

	log.Info("starting email API server",
		"port", cfg.Port,
		"environment", cfg.Environment,
	)

	// Initialize Gmail sender
	gmailSender, err := email.NewGmailSender()
	if err != nil {
		log.Error("failed to initialize Gmail sender", "error", err)
		os.Exit(1)
	}
	log.Info("Gmail sender initialized successfully")

	// Initialize email API router
	router := emailapi.NewEmailAPIRouter(gmailSender, log)

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
		log.Info("email API server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down email API server")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown error", "error", err)
		os.Exit(1)
	}

	log.Info("email API server stopped")
}

