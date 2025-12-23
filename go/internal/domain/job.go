package domain

import "time"

// Job represents a pipeline execution job
type Job struct {
	ID          string                 `json:"id"`
	BusinessID  string                 `json:"businessId"`
	PipelineKey string                 `json:"pipelineKey"`
	Status      string                 `json:"status"` // "pending" | "running" | "completed" | "failed"
	Steps       []JobStep              `json:"steps"`
	Input       map[string]any         `json:"input"`
	Result      map[string]any         `json:"result"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// JobStep represents a single step in a job
type JobStep struct {
	Name     string         `json:"name"`
	Status   string         `json:"status"` // "ok" | "skipped" | "failed"
	Critical bool           `json:"critical"`
	Error    *string        `json:"error,omitempty"`
	Details  map[string]any `json:"details,omitempty"`
}

