package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/util"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// EmailAnalysisStatsHandler handles email analysis statistics
type EmailAnalysisStatsHandler struct {
	logger *slog.Logger
}

// NewEmailAnalysisStatsHandler creates a new stats handler
func NewEmailAnalysisStatsHandler(logger *slog.Logger) *EmailAnalysisStatsHandler {
	return &EmailAnalysisStatsHandler{
		logger: logger,
	}
}

// StatsResponse represents the stats data
type StatsResponse struct {
	TotalProcessed    int                    `json:"total_processed"`
	TotalSkipped      int                    `json:"total_skipped"`
	CurrentJobID      string                 `json:"current_job_id"`
	CurrentJobName    string                 `json:"current_job_name"`
	LastRun           time.Time              `json:"last_run"`
	ProcessedIDsCount int                    `json:"processed_ids_count"`
	Jobs              []JobStats             `json:"jobs"`
	ActiveAgents      []AgentInfo            `json:"active_agents"`
	SpreadsheetID     string                 `json:"spreadsheet_id"`
	SpreadsheetURL    string                 `json:"spreadsheet_url"`
	LastUpdated       time.Time              `json:"last_updated"`
}

// JobStats represents job statistics
type JobStats struct {
	JobID        string    `json:"job_id"`
	JobName      string    `json:"job_name"`
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	TotalProcessed int     `json:"total_processed"`
	TotalSkipped   int     `json:"total_skipped"`
	Status        string   `json:"status"`
	AgentIDs      []string `json:"agent_ids"`
}

// AgentInfo represents active agent information
type AgentInfo struct {
	AgentID   string    `json:"agent_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// HandleStats returns current analysis statistics
func (h *EmailAnalysisStatsHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	spreadsheetID := r.URL.Query().Get("spreadsheet_id")
	if spreadsheetID == "" {
		util.WriteError(w, http.StatusBadRequest, "spreadsheet_id parameter required")
		return
	}

	ctx := r.Context()
	stats, err := h.getStats(ctx, spreadsheetID)
	if err != nil {
		h.logger.Error("failed to get stats", "error", err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, stats)
}

// getStats retrieves statistics from the spreadsheet
func (h *EmailAnalysisStatsHandler) getStats(ctx context.Context, spreadsheetID string) (*StatsResponse, error) {
	service, err := initSheetsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init sheets: %w", err)
	}

	stats := &StatsResponse{
		SpreadsheetID:  spreadsheetID,
		SpreadsheetURL: fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", spreadsheetID),
		LastUpdated:    time.Now(),
	}

	// Get state
	stateResp, err := service.Spreadsheets.Values.Get(spreadsheetID, "State!A2:F2").Context(ctx).Do()
	if err == nil && len(stateResp.Values) > 0 {
		row := stateResp.Values[0]
		if len(row) > 1 {
			if total, err := strconv.Atoi(fmt.Sprintf("%v", row[1])); err == nil {
				stats.TotalProcessed = total
			}
		}
		if len(row) > 2 {
			if lastRun, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", row[2])); err == nil {
				stats.LastRun = lastRun
			}
		}
		if len(row) > 3 {
			if count, err := strconv.Atoi(fmt.Sprintf("%v", row[3])); err == nil {
				stats.ProcessedIDsCount = count
			}
		}
		if len(row) > 4 {
			stats.CurrentJobID = fmt.Sprintf("%v", row[4])
		}
		if len(row) > 5 {
			stats.CurrentJobName = fmt.Sprintf("%v", row[5])
		}
	}

	// Get job stats
	jobResp, err := service.Spreadsheets.Values.Get(spreadsheetID, "Job Stats!A2:H").Context(ctx).Do()
	if err == nil && len(jobResp.Values) > 0 {
		jobs := make([]JobStats, 0, len(jobResp.Values))
		for _, row := range jobResp.Values {
			if len(row) < 6 {
				continue
			}
			
			job := JobStats{
				JobID: fmt.Sprintf("%v", row[0]),
				JobName: fmt.Sprintf("%v", row[1]),
				Status: fmt.Sprintf("%v", row[6]),
			}
			
			if len(row) > 2 {
				if started, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", row[2])); err == nil {
					job.StartedAt = started
				}
			}
			if len(row) > 3 && row[3] != nil && fmt.Sprintf("%v", row[3]) != "" {
				if completed, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", row[3])); err == nil {
					job.CompletedAt = &completed
				}
			}
			if len(row) > 4 {
				if processed, err := strconv.Atoi(fmt.Sprintf("%v", row[4])); err == nil {
					job.TotalProcessed = processed
				}
			}
			if len(row) > 5 {
				if skipped, err := strconv.Atoi(fmt.Sprintf("%v", row[5])); err == nil {
					job.TotalSkipped = skipped
				}
			}
			if len(row) > 7 {
				agentIDsStr := fmt.Sprintf("%v", row[7])
				if agentIDsStr != "" {
					// Parse comma-separated agent IDs
					parts := strings.Split(agentIDsStr, ",")
					job.AgentIDs = make([]string, 0, len(parts))
					for _, part := range parts {
						part = strings.TrimSpace(part)
						if part != "" {
							job.AgentIDs = append(job.AgentIDs, part)
						}
					}
				}
			}
			
			jobs = append(jobs, job)
		}
		stats.Jobs = jobs
	}

	// Get active locks (agents)
	lockResp, err := service.Spreadsheets.Values.Get(spreadsheetID, "Locks!A2:D").Context(ctx).Do()
	if err == nil && len(lockResp.Values) > 0 {
		now := time.Now()
		agents := make([]AgentInfo, 0)
		
		for _, row := range lockResp.Values {
			if len(row) < 4 {
				continue
			}
			
			status := fmt.Sprintf("%v", row[3])
			if status != "ACTIVE" {
				continue
			}
			
			expiresStr := fmt.Sprintf("%v", row[2])
			expiresAt, err := time.Parse(time.RFC3339, expiresStr)
			if err != nil {
				continue
			}
			
			// Only include non-expired locks
			if expiresAt.After(now) {
				createdAt := time.Now()
				if len(row) > 1 {
					if created, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", row[1])); err == nil {
						createdAt = created
					}
				}
				
				agents = append(agents, AgentInfo{
					AgentID:   fmt.Sprintf("%v", row[0]),
					ExpiresAt: expiresAt,
					CreatedAt: createdAt,
				})
			}
		}
		stats.ActiveAgents = agents
	}

	// Calculate total skipped from jobs
	for _, job := range stats.Jobs {
		stats.TotalSkipped += job.TotalSkipped
	}

	return stats, nil
}

// initSheetsService initializes Google Sheets service
func initSheetsService(ctx context.Context) (*sheets.Service, error) {
	credsJSON := os.Getenv("GMAIL_CREDENTIALS_JSON")
	if credsJSON == "" {
		return nil, fmt.Errorf("GMAIL_CREDENTIALS_JSON not set")
	}

	var credsData []byte
	if _, err := os.Stat(credsJSON); err == nil {
		credsData, _ = os.ReadFile(credsJSON)
	} else {
		credsData = []byte(credsJSON)
	}

	config, err := google.JWTConfigFromJSON(credsData, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	// Domain-wide delegation: impersonate a real user.
	// Without Subject, reads can silently return empty/permission issues -> dashboard shows zeros.
	if subj := os.Getenv("GMAIL_FROM"); subj != "" {
		config.Subject = subj
	} else {
		// Reasonable default for this repo's primary use case.
		config.Subject = "team@stlpartyhelpers.com"
	}

	client := config.Client(ctx)
	return sheets.NewService(ctx, option.WithHTTPClient(client))
}
