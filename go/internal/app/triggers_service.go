package app

import (
	"context"
	"fmt"

	"github.com/bizops360/go-api/internal/config"
	"github.com/bizops360/go-api/internal/domain"
	"github.com/bizops360/go-api/internal/ports"
)

// TriggersService handles trigger-based pipeline execution
type TriggersService struct {
	businessLoader *config.BusinessLoader
	pipelineRunner *domain.PipelineRunner
	jobsRepo       ports.JobsRepo
}

// NewTriggersService creates a new triggers service
func NewTriggersService(
	businessLoader *config.BusinessLoader,
	pipelineRunner *domain.PipelineRunner,
	jobsRepo ports.JobsRepo,
) *TriggersService {
	return &TriggersService{
		businessLoader: businessLoader,
		pipelineRunner: pipelineRunner,
		jobsRepo:       jobsRepo,
	}
}

// Run executes a trigger pipeline
func (s *TriggersService) Run(ctx context.Context, req *TriggerRequest) (*domain.PipelineResult, error) {
	// Load business config
	business, err := s.businessLoader.LoadBusiness(ctx, req.BusinessID)
	if err != nil {
		return nil, fmt.Errorf("failed to load business: %w", err)
	}

	// Resolve pipeline key from trigger key
	pipelineKey := req.PipelineKey
	if pipelineKey == "" && req.TriggerKey != "" {
		var ok bool
		pipelineKey, ok = business.Pipelines.Triggers[req.TriggerKey]
		if !ok {
			return nil, fmt.Errorf("trigger key '%s' not found in business config", req.TriggerKey)
		}
	}
	if pipelineKey == "" {
		return nil, fmt.Errorf("pipeline key not specified and trigger key not found")
	}

	// Load pipeline
	pipeline, err := s.businessLoader.LoadPipeline(ctx, pipelineKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load pipeline: %w", err)
	}

	// Build pipeline context
	pctx := &domain.PipelineContext{
		BusinessID:  req.BusinessID,
		PipelineKey: pipelineKey,
		Source:      req.Source,
		DryRun:      req.DryRun,
		Fields:      req.Payload,
		Options:     make(map[string]any),
		Resource:    req.Resource,
		RequestID:   req.RequestID,
	}

	// Run pipeline
	result, job := s.pipelineRunner.Run(ctx, pipeline, pctx)

	// Save job
	if job != nil {
		if err := s.jobsRepo.Save(ctx, job); err != nil {
			// Log error but don't fail the request
			_ = err
		}
	}

	return result, nil
}

// TriggerRequest represents a trigger request
type TriggerRequest struct {
	BusinessID  string
	TriggerKey  string
	PipelineKey string
	Source      string
	Resource    *domain.ResourceContext
	Payload     map[string]any
	DryRun      bool
	RequestID   string
}

