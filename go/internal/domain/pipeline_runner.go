package domain

import (
	"context"
	"fmt"
)

// Action is the interface that all pipeline actions must implement
type Action interface {
	Name() string
	Execute(ctx context.Context, pctx *PipelineContext) JobStep
}

// PipelineRunner executes pipelines by running actions in sequence
type PipelineRunner struct {
	actions map[string]Action
}

// NewPipelineRunner creates a new pipeline runner with registered actions
func NewPipelineRunner(actions map[string]Action) *PipelineRunner {
	return &PipelineRunner{
		actions: actions,
	}
}

// Run executes a pipeline definition with the given context
func (pr *PipelineRunner) Run(ctx context.Context, pipeline *PipelineDefinition, pctx *PipelineContext) (*PipelineResult, *Job) {
	job := &Job{
		ID:          pctx.RequestID,
		BusinessID:  pctx.BusinessID,
		PipelineKey: pctx.PipelineKey,
		Status:      "running",
		Steps:       []JobStep{},
		Input:       pctx.Fields,
		Result:      make(map[string]any),
	}

	result := &PipelineResult{
		PipelineKey: pipeline.Key,
		BusinessID:  pctx.BusinessID,
		DryRun:      pctx.DryRun,
		Steps:       []JobStep{},
		JobID:       job.ID,
	}

	// Execute each action in sequence
	for _, actionDef := range pipeline.Actions {
		action, exists := pr.actions[actionDef.Name]
		if !exists {
			step := JobStep{
				Name:     actionDef.Name,
				Status:   "failed",
				Critical: actionDef.Critical,
				Error:    stringPtr(fmt.Sprintf("action '%s' not found", actionDef.Name)),
			}
			job.Steps = append(job.Steps, step)
			result.Steps = append(result.Steps, step)
			
			if actionDef.Critical {
				result.Success = false
				result.Error = stringPtr(fmt.Sprintf("critical action '%s' failed: action not found", actionDef.Name))
				job.Status = "failed"
				return result, job
			}
			continue
		}

		// Execute the action
		step := action.Execute(ctx, pctx)
		job.Steps = append(job.Steps, step)
		result.Steps = append(result.Steps, step)

		// Check if step failed
		if step.Status == "failed" {
			if step.Critical {
				result.Success = false
				result.Error = step.Error
				job.Status = "failed"
				return result, job
			}
			// Non-critical failure, continue
		}
	}

	// All steps completed successfully
	result.Success = true
	job.Status = "completed"
	return result, job
}

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}

