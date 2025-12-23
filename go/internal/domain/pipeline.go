package domain

// PipelineDefinition represents a pipeline configuration
type PipelineDefinition struct {
	Key         string             `yaml:"key" json:"key"`
	Description string             `yaml:"description" json:"description"`
	Actions     []ActionDefinition `yaml:"actions" json:"actions"`
}

// ActionDefinition defines an action within a pipeline
type ActionDefinition struct {
	Name     string         `yaml:"name" json:"name"`
	Critical bool           `yaml:"critical" json:"critical"`
	Config   map[string]any `yaml:"config" json:"config"`
}

// PipelineContext holds context for pipeline execution
type PipelineContext struct {
	BusinessID  string
	PipelineKey string
	Source      string
	DryRun      bool
	Fields      map[string]any
	Options     map[string]any
	Resource    *ResourceContext
	RequestID   string
}

// ResourceContext holds information about the triggering resource
type ResourceContext struct {
	Type    string         `json:"type"`
	BoardID *int64         `json:"boardId,omitempty"`
	ItemID  *int64         `json:"itemId,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

// PipelineResult represents the result of pipeline execution
type PipelineResult struct {
	Success     bool       `json:"success"`
	PipelineKey string     `json:"pipelineKey"`
	BusinessID  string     `json:"businessId"`
	DryRun      bool       `json:"dryRun"`
	Steps       []JobStep  `json:"steps"`
	JobID       string     `json:"jobId,omitempty"`
	Error       *string    `json:"error,omitempty"`
}

