package handlers

import (
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/bizops360/go-api/internal/util"
)

// Commit represents a git commit
type Commit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

// CommitsHandler handles git commits endpoint
type CommitsHandler struct{}

// NewCommitsHandler creates a new commits handler
func NewCommitsHandler() *CommitsHandler {
	return &CommitsHandler{}
}

// HandleCommits handles GET /api/commits
func (h *CommitsHandler) HandleCommits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Get limit from query parameter (default 10)
	limit := "10"
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		limit = limitParam
	}

	// Run git log command
	cmd := exec.Command("git", "log", "--pretty=format:%H|%s|%an|%ad", "--date=iso", "-n", limit)
	output, err := cmd.Output()
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, "failed to fetch commits: "+err.Error())
		return
	}

	// Parse git log output
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	commits := make([]Commit, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 4 {
			continue
		}

		hash := parts[0]
		message := parts[1]
		author := parts[2]
		dateStr := parts[3]

		// Parse and format date
		date, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		if err != nil {
			// Try alternative format
			date, err = time.Parse(time.RFC3339, dateStr)
			if err != nil {
				date = time.Now()
			}
		}
		dateFormatted := date.Format("Jan 2, 2006 3:04 PM")

		commits = append(commits, Commit{
			Hash:    hash[:8], // Short hash
			Message: message,
			Author:  author,
			Date:    dateFormatted,
		})
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"commits": commits,
		"count":   len(commits),
	})
}
